package main

import (
	"../proxy/container"
	"../proxy/models"
	"../proxy/shared"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/robfig/cron.v2"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//"github.com/jinzhu/gorm"
)

var sharedContext *shared.SharedContext

type Job struct {
	CronJob models.CronJob
	EntryID cron.EntryID
}

type FuncJob struct {
	Account models.Account
	CronJob models.CronJob
}

func (f FuncJob) Run() {
	defer func() {
		err := recover()

		if err != nil {
			log.Printf(fmt.Sprintf("[cj %d] panicked with: %s"), f.CronJob.Id)
		}
	}()

	log.Printf("[cj %d] running func job", f.CronJob.Id)

	cjLog := models.CronJobLog{
		CronJobId: f.CronJob.Id,
		Added_at:  time.Now(),
	}
	sharedContext.PersistentDB.Save(&cjLog)

	out, err := exec.Command("id", "-u", f.Account.Name).Output()
	if err != nil {
		log.Printf("[cj %d] cannot get account user %s uid", f.CronJob.Id, f.Account.Name)
		return
	}

	uid := strings.Replace(string(out), "\n", "", 1)

	cmd := strings.Split(f.CronJob.Command, " ")

	s := container.Shell{
		LogPrefix:     "[cron]",
		AccountName:   f.Account.Name,
		Cmd:           cmd,
		Tty:           false,
		AccountUid:    uid,
		SharedContext: sharedContext,
	}

	endpoint := "unix:///var/run/docker.sock"
	s.DockerClient, err = docker.NewClient(endpoint)
	if err != nil {
		log.Printf("[cj %d] error while creating docker client: %s", f.CronJob.Id, err)
		return
	}

	s.GetDockerImages()

	s.Environment = strings.Replace(f.CronJob.Image, "-base-shell", "", 1)

	shell_image, err := s.BuildShellImage(s.Environment)
	if err != nil {
		s.LogError(err)
		return
	}

	s.WorkingDir = f.CronJob.Directory
	container, err := s.CreateContainer(shell_image)
	if err != nil {
		s.LogError(fmt.Errorf("couldn't create container %s (image: %s)", err, shell_image))
		return
	}

	defer func() {
		s.Log("info", "cleanup shell %s %s", s.AccountName, shell_image)

		if container != nil {
			err := s.RemoveContainer()

			if err != nil {
				s.LogError(fmt.Errorf("error while cleaning up %s", err))
			}
		}
	}()

	err = s.StartContainer()
	if err != nil {
		s.LogError(fmt.Errorf("cannot start container ", err))
		return
	}

	buf := bytes.NewBuffer(nil)

	errs := make(chan error)
	go func() {
		errs <- s.DockerClient.AttachToContainer(docker.AttachToContainerOptions{
			Container:    container.ID,
			OutputStream: buf,
			ErrorStream:  buf,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			RawTerminal:  false,
		})
	}()
	if err != nil {
		s.LogError(fmt.Errorf("cannot attach to container ", err))
		return
	}
	myerr := <-errs

	if myerr != nil {
		s.LogError(fmt.Errorf("attach error %s", err))
		return
	}

	code, err := s.DockerClient.WaitContainer(container.ID)
	if err != nil {
		s.LogError(fmt.Errorf("cannot wait for the container ", err))
		return
	}

	var lines = ""
	var line = ""
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line = scanner.Text()
		lines = fmt.Sprintf("%s\n%s", lines, line)
	}

	if err := scanner.Err(); err != nil {
		s.LogError(errors.New(fmt.Sprintf("error encountered while reading output: %s", err)))
		return
	}

	cjLog.Log = lines
	cjLog.Duration = time.Now().Sub(cjLog.Added_at).Seconds()
	cjLog.ExitCode = code

	if code == 0 {
		cjLog.Success = true
	}

	sharedContext.PersistentDB.Save(&cjLog)

	var cj models.CronJob
	err = sharedContext.PersistentDB.Where("id = ?", strconv.Itoa(f.CronJob.Id)).First(&cj).Error
	if err == nil {
		cj.Success = cjLog.Success
		sharedContext.PersistentDB.Save(cj)
	}

	s.Log("info", "[cj %d] session (exec, %s (%s)) closed", string(f.CronJob.Id), s.AccountName, s.AccountUid)
}

func main() {
	sharedContext = &shared.SharedContext{}
	sharedContext.OpenDB("../manager/db.sqlite3")

	jobs := make(map[int]Job)

	c := cron.New()
	c.Start()

	prev, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 1 15:04:05 -0700 MST 2000")

	if err != nil {
		log.Fatalf("error parsing date: %s", err)
	}

	for {
		var cronjobs []models.CronJob
		sharedContext.PersistentDB.Where("updated_at > ?", prev).Find(&cronjobs)

		for _, cj := range cronjobs {
			if cj.UpdatedAt.After(prev) {
				prev = cj.UpdatedAt
			}

			if v, ok := jobs[cj.Id]; ok {
				//remove old job if exists
				log.Printf("removing job id %d (entry id: %d)", cj.Id, v.EntryID)
				c.Remove(v.EntryID)
				delete(jobs, cj.Id)
			}

			//get account info
			var account models.Account
			if err := sharedContext.PersistentDB.Where("id = ?", cj.Account_id).First(&account).Error; err != nil {
				log.Printf("cannot find account with id %d for cj %d", cj.Account_id, cj.Id)
				continue
			}

			//add new job
			if !cj.Enabled {
				log.Printf("cron job disabled (%d)", cj.Id)
				continue
			}

			f := FuncJob{
				CronJob: cj,
				Account: account,
			}
			id, err := c.AddJob(cj.Cron_expression, f)
			if err != nil {
				log.Printf("error adding job id %d with spec: %s", cj.Id, cj.Cron_expression)
				continue
			}

			log.Printf("adding job id %d with spec: %s as entry id: %d", cj.Id, cj.Cron_expression, id)

			jobs[cj.Id] = Job{
				CronJob: cj,
				EntryID: id,
			}
		}

		time.Sleep(time.Second * 1)
	}
}
