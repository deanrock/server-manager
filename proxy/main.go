package main

import (
	"./container"
	"./controllers"
	"./models"
	"./realtime"
	"./shared"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samalba/dockerclient"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

var sharedContext *shared.SharedContext
var mydocker *dockerclient.DockerClient

type MyContainer struct {
	AppId       int    `json:"app_id"`
	AppName     string `json:"app_name"`
	AccountName string `json:"account_name"`
	Up          bool   `json:"up"`
	Status      string `json:"status"`
	ImageName   string `json:"image_name"`
	Memory      int    `json:"memory"`
}

func Containers(c *gin.Context) {
	containers, err := container.GetAllContainers(sharedContext)
	if err != nil {
		c.AbortWithStatus(500)
	}

	var apps []models.App
	sharedContext.PersistentDB.Find(&apps)

	var accounts []models.Account
	sharedContext.PersistentDB.Find(&accounts)

	userAccess := c.MustGet("userAccess").([]models.UserAccess)

	var allowedContainers []MyContainer

	var images []models.Image
	sharedContext.PersistentDB.Find(&images)

	for _, app := range apps {
		var account models.Account
		found := false
		for _, acc := range accounts {
			if acc.Id == app.Account_id {
				for _, ua := range userAccess {
					fmt.Println(ua)
					if ua.Account_id == acc.Id && ua.AppAccess == true {
						account = acc
						found = true
						break
					}
				}
				break
			}
		}

		if found {
			name := fmt.Sprintf("/app-%s-%s", account.Name, app.Name)
			var con docker.APIContainers
			found := false

			for _, c := range containers {
				for _, n := range c.Names {
					if n == name {
						con = c
						found = true
						break
					}
				}
			}

			if found {
				up := false

				if strings.Contains(con.Status, "Up ") {
					up = true
				}

				imageName := ""
				for _, i := range images {
					if i.Id == app.Image_id {
						imageName = i.Name
						break
					}
				}

				allowedContainers = append(allowedContainers, MyContainer{
					AppId:       app.Id,
					AppName:     app.Name,
					AccountName: account.Name,
					Up:          up,
					Status:      con.Status,
					ImageName:   imageName,
					Memory:      app.Memory,
				})
			}
		}
	}

	c.JSON(200, allowedContainers)
}

func containerLogsHandler(c *gin.Context) {
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	options := dockerclient.LogOptions{
		Follow:     true,
		Stdout:     true,
		Stderr:     true,
		Timestamps: true,
	}

	a := models.AccountFromContext(c)

	id, e := strconv.ParseInt(c.Params.ByName("id"), 0, 64)
	if e != nil {
		c.String(400, "")
		return
	}

	var container docker.APIContainers
	var found = false

	for _, app := range a.Apps() {
		if app.Id == int(id) {
			log.Println(app.Id)
			log.Println(app)

			containers, _ := sharedContext.DockerClient.ListContainers(docker.ListContainersOptions{})

			name := fmt.Sprintf("/%s", app.ContainerName(a.Name))
			for _, c := range containers {
				for _, n := range c.Names {
					if n == name {
						container = c
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		http.Error(c.Writer, "", 404)
		return
	}

	reader, err := mydocker.ContainerLogs(container.ID, &options)
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Printf("Read Error:", err)
			return
		}

		if len(str) > 8 {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(str[8:])); err != nil {
				conn.Close()
				break
			}
		}
	}
}

func proxyRequest(c *gin.Context) {
	director := func(req *http.Request) {
		req = c.Request
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:5555"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := readCookie(c)

		if userId != nil {
			u, err := models.FindUserById(sharedContext, *userId)

			if err != nil {
				c.AbortWithError(500, errors.New("cannot handle"))
				c.Abort()
			}

			c.Set("uid", userId)
			c.Set("user", *u)

			var access []models.UserAccess
			sharedContext.PersistentDB.Where("user_id = ?", userId).Find(&access)
			c.Set("userAccess", access)
		} else {
			//c.AbortWithError(401, errors.New("Unauthorized"))
			c.Redirect(http.StatusSeeOther, "/accounts/login/?next=/")
			c.Abort()
		}

		c.Next()
	}
}

func RequireAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Params.ByName("name")
		account := models.GetAccountByName(name, sharedContext)

		if account != nil {
			userAccess := c.MustGet("userAccess").([]models.UserAccess)

			found := false
			for _, ua := range userAccess {
				if ua.Account_id == account.Id {
					found = true
				}
			}

			if found {
				c.Set("account", account)
			} else {
				c.AbortWithError(401, errors.New("unauthorized access to account"))
			}
		} else {
			c.AbortWithError(404, errors.New("Not Found"))
		}

		c.Next()
	}
}

func RequireUserAccess(accessName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Params.ByName("name")
		account := models.GetAccountByName(name, sharedContext)

		if account != nil {
			userAccess := c.MustGet("userAccess").([]models.UserAccess)

			found := false
			for _, ua := range userAccess {
				if ua.Account_id == account.Id {
					if accessName == "shell_access" && ua.ShellAccess == true {
						found = true
					} else if accessName == "app_access" && ua.AppAccess == true {
						found = true
					} else if accessName == "cronjob_access" && ua.CronjobAccess == true {
						found = true
					} else if accessName == "domain_access" && ua.DomainAccess == true {
						found = true
					} else if accessName == "ssh_access" && ua.SshAccess == true {
						found = true
					} else if accessName == "database_access" && ua.DatabaseAccess == true {
						found = true
					}
				}
			}

			if !found {
				c.AbortWithError(401, errors.New("user doesnt have access to this resource"))
			}
		} else {
			c.AbortWithError(404, errors.New("Not Found"))
		}

		c.Next()
	}
}

func RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)

		if user.Is_staff {
			c.Next()
		} else {
			c.AbortWithError(401, errors.New("Unauthorized"))
			c.Abort()
		}
	}
}

type Config struct {
	Server_name             string `json:"server_name"`
	Mysql_connection_string string `json:"mysql_connection_string"`
}

type Profile struct {
	Server_name string      `json:"server_name"`
	User        models.User `json:"user"`
}

func main() {
	c, _ := ioutil.ReadFile("../config.json")
	dec := json.NewDecoder(bytes.NewReader(c))
	var config Config
	dec.Decode(&config)

	//Docker
	// Init the client
	mydocker, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

	sharedContext = &shared.SharedContext{}
	sharedContext.OpenDB("../manager/db.sqlite3")
	sharedContext.PersistentDB.LogMode(true)
	sharedContext.PersistentDB.AutoMigrate(&models.CronJob{})
	sharedContext.PersistentDB.AutoMigrate(&models.CronJobLog{})
	sharedContext.PersistentDB.AutoMigrate(&models.UserAccess{})
	sharedContext.PersistentDB.Model(&models.UserAccess{}).AddUniqueIndex("idx_user_account", "user_id", "account_id")
	sharedContext.PersistentDB.AutoMigrate(&models.Task{})
	sharedContext.PersistentDB.AutoMigrate(&models.TaskLog{})
	sharedContext.PersistentDB.AutoMigrate(&models.SSHPassword{})
	sharedContext.WebsocketHandler = realtime.NewWebsocketHandler()

	//go-dockerclient
	endpoint := "unix:///var/run/docker.sock"
	sharedContext.DockerClient, _ = docker.NewClient(endpoint)

	//cancell all tasks after restart
	models.CancelAllTasks(sharedContext)

	//HTTP
	r := gin.New()
	r.Use(gin.Logger())
	///////r.Use(gin.Recovery()) !!! DONT USE gin.Recovery or gin.Default, because it ignores middleware if it panics ...
	//                               which means that panic in Authentication() or RequireStaff() can allow unauthorized users access to all paths

	r.LoadHTMLGlob("templates/*")

	authorized := r.Group("/")

	authorized.Use(Authentication())
	{
		authorized.GET("/ws/", func(c *gin.Context) {
			user := c.MustGet("user").(models.User)
			sharedContext.WebsocketHandler.WsHandler(c, user.Id, user.Is_staff, func(uid int) {
				//running tasks
				msg, _ := json.Marshal(struct {
					Type  string        `json:"type"`
					Tasks []models.Task `json:"tasks"`
				}{
					Type:  "my-running-tasks",
					Tasks: models.RunningTasksForUser(sharedContext, user),
				})

				sharedContext.WebsocketHandler.SendToUser(msg, uid)
			})
		})

		//profile
		authorized.GET("/api/v1/profile", func(c *gin.Context) {
			uid := c.MustGet("uid").(*int)
			user, _ := models.FindUserById(sharedContext, *uid)
			p := Profile{
				Server_name: config.Server_name,
				User:        *user,
			}

			c.JSON(200, p)
		})

		userSSHKeys := &controllers.UserSSHKeysAPI{
			Context: sharedContext,
		}

		authorized.GET("/api/v1/profile/ssh-keys", userSSHKeys.ListKeys)
		authorized.GET("/api/v1/profile/ssh-keys/:id", userSSHKeys.GetKey)
		authorized.PUT("/api/v1/profile/ssh-keys/:id", userSSHKeys.EditKey)
		authorized.DELETE("/api/v1/profile/ssh-keys/:id", userSSHKeys.DeleteKey)
		authorized.POST("/api/v1/profile/ssh-keys", userSSHKeys.EditKey)

		authorized.GET("/api/v1/shells", func(c *gin.Context) {
			s := container.Shell{}
			s.GetDockerImages()
			c.JSON(200, s.ShellImages)
		})

		//images
		authorized.GET("/api/v1/images", func(c *gin.Context) {
			images := models.GetImages(sharedContext)

			c.JSON(200, images)
		})

		//accounts
		accounts := &controllers.AccountsAPI{
			Context: sharedContext,
		}

		authorized.GET("/api/v1/accounts", accounts.ListAccounts)
		authorized.GET("/api/v1/all-accounts", RequireStaff(), accounts.ListAllAccounts)
		authorized.GET("/api/v1/containers", Containers)

		requiresAccount := authorized.Group("/api/v1/accounts/:name")

		requiresAccount.Use(RequireAccount())
		{
			requiresAccount.GET("", accounts.GetAccountByName)

			//apps
			apps := &controllers.AppsAPI{
				Context: sharedContext,
			}

			requiresAccount.GET("/apps", RequireUserAccess("app_access"), apps.ListApps)
			requiresAccount.GET("/apps/:id", RequireUserAccess("app_access"), apps.GetApp)
			requiresAccount.PUT("/apps/:id", RequireUserAccess("app_access"), apps.EditApp)
			requiresAccount.POST("/apps", RequireUserAccess("app_access"), apps.EditApp)
			requiresAccount.POST("/apps/:id/redeploy", RequireUserAccess("app_access"), apps.RedeployApp)
			requiresAccount.POST("/apps/:id/start", RequireUserAccess("app_access"), apps.StartApp)
			requiresAccount.POST("/apps/:id/stop", RequireUserAccess("app_access"), apps.StopApp)
			requiresAccount.GET("/apps/:id/logs", RequireUserAccess("app_access"), containerLogsHandler)

			//shell
			requiresAccount.GET("/shell", RequireUserAccess("shell_access"), func(c *gin.Context) {
				container.WebSocketShell(c, sharedContext)
			})

			//cronjobs
			cronJobs := &controllers.CronJobsAPI{
				Context: sharedContext,
			}

			requiresAccount.GET("/cronjobs", RequireUserAccess("cronjob_access"), cronJobs.ListCronjobs)
			requiresAccount.GET("/cronjobs/:id", RequireUserAccess("cronjob_access"), cronJobs.GetCronjob)
			requiresAccount.GET("/cronjobs/:id/log", RequireUserAccess("cronjob_access"), cronJobs.GetCronjobLog)
			requiresAccount.PUT("/cronjobs/:id", RequireUserAccess("cronjob_access"), cronJobs.EditCronjob)
			requiresAccount.POST("/cronjobs", RequireUserAccess("cronjob_access"), cronJobs.AddCronjob)

			//domains
			domains := &controllers.DomainsAPI{
				Context: sharedContext,
			}

			requiresAccount.GET("/domains", RequireUserAccess("domain_access"), domains.ListDomains)
			requiresAccount.GET("/domains/:id", RequireUserAccess("domain_access"), domains.GetDomain)
			requiresAccount.PUT("/domains/:id", RequireUserAccess("domain_access"), domains.EditDomain)
			requiresAccount.DELETE("/domains/:id", RequireUserAccess("domain_access"), domains.DeleteDomain)
			requiresAccount.POST("/domains", RequireUserAccess("domain_access"), domains.EditDomain)
			requiresAccount.POST("/domains/sync", RequireUserAccess("domain_access"), domains.SyncDomains)

			//databases
			databases := &controllers.DatabasesAPI{
				Context: sharedContext,
			}

			requiresAccount.GET("/databases", RequireUserAccess("database_access"), databases.ListDatabases)
			requiresAccount.GET("/databases/:id", RequireUserAccess("database_access"), databases.GetDatabase)
			requiresAccount.PUT("/databases/:id", RequireUserAccess("database_access"), databases.EditDatabase)
			requiresAccount.POST("/databases", RequireUserAccess("database_access"), databases.EditDatabase)

			//ssh passwords
			sshPasswords := &controllers.SSHPasswordsAPI{
				Context: sharedContext,
			}

			requiresAccount.GET("/ssh-passwords", RequireStaff(), sshPasswords.ListPasswords)
			requiresAccount.DELETE("/ssh-passwords/:id", RequireStaff(), sshPasswords.DeletePassword)
			requiresAccount.POST("/ssh-passwords", RequireStaff(), sshPasswords.AddPassword)
		}

		//users
		users := &controllers.UsersAPI{
			Context: sharedContext,
		}

		authorized.GET("/api/v1/profile/access/:account", users.GetMyAccess)

		u := authorized.Group("/api/v1/users", RequireStaff())
		{
			u.GET("", users.ListUsers)
			u.GET(":id", users.GetUser)
			u.GET(":id/access", users.GetAccess)
			u.POST(":id/access/:account", users.SetAccess)
			u.DELETE(":id/access/:account", users.RemoveAccess)
		}

		//sync
		sync := &controllers.SyncAPI{
			Context: sharedContext,
		}

		authorized.GET("/api/v1/sync/images", RequireStaff(), sync.GetImages)
		authorized.POST("/api/v1/sync/images/:name", RequireStaff(), sync.SyncImage)
		authorized.POST("/api/v1/sync/web-servers", RequireStaff(), sync.SyncWebServers)

		//tasks
		tasks := &controllers.TasksAPI{
			Context: sharedContext,
		}

		authorized.GET("/api/v1/tasks", tasks.ListMyTasks)
		authorized.GET("/api/v1/tasks/:id", tasks.GetTask)
		authorized.GET("/api/v1/tasks/:id/log", tasks.GetTaskLog)

		//angular
		authorized.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/a/*params", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/sync/*params", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/containers", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/users/*params", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/profile/*params", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})

		authorized.GET("/tasks/*params", func(c *gin.Context) {
			c.HTML(200, "index.tmpl", nil)
		})
	}

	s := &http.Server{
		Addr:           ":4444",
		Handler:        r,
		ReadTimeout:    60 * time.Minute,
		WriteTimeout:   60 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	r.Static("/static/ace", "../static/ace")
	r.Static("/static/bootstrap", "../static/bootstrap")
	r.Static("/static/css", "../static/css")
	r.Static("/static/js", "../static/js")
	r.Static("/static/images", "../static/images")
	r.Static("/static/vendor", "../static/vendor")
	r.Static("/static/templates", "../static/templates")
	r.Static("/static/admin", "../env/lib/python2.7/site-packages/django/contrib/admin/static/admin")

	r.NoRoute(proxyRequest)

	//certFile := "../ssl.crt"
	//keyFile := "../ssl.key"

	/*if _, err := os.Stat(certFile); err == nil {
	    log.Fatal(s.ListenAndServeTLS(certFile, keyFile))
	}else{*/
	log.Fatal(s.ListenAndServe())
	//}
}
