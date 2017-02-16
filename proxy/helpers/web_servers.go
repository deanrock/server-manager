package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"../container"
	"../models"
	"../shared"
	"github.com/fsouza/go-dockerclient"
)

func removeConfigFiles(path string, prefix string) error {
	m, err := filepath.Glob(fmt.Sprintf("%s/*", path))
	if err != nil {
		return err
	}

	prefix = fmt.Sprintf("%s/%s", path, prefix)

	for _, file := range m {
		if strings.HasPrefix(file, prefix) {
			os.Remove(file)
		}
	}

	return nil
}

func SyncWebServers(sharedContext *shared.SharedContext, task models.Task, account *models.Account) bool {
	success := true

	// get all cointainers
	containers, err := container.GetAllContainers(sharedContext)
	if err != nil {
		return false
	}

	// get images
	images := models.GetImages(sharedContext)

	var accounts []models.Account

	prefix := ""
	if account == nil {
		// get all accounts
		sharedContext.PersistentDB.Find(&accounts)
	} else {
		prefix = fmt.Sprintf("%s_", account.Name)

		// add account to array
		accounts = append(accounts, *account)
	}

	// remove config for all accounts
	err = removeConfigFiles("/etc/nginx/manager", prefix)
	if err != nil {
		task.Log(fmt.Sprintf("error while deleting nginx config files: %s", err), "error", sharedContext)
	}

	err = removeConfigFiles("/etc/apache2/manager", prefix)
	if err != nil {
		task.Log(fmt.Sprintf("error while deleting apache config files: %s", err), "error", sharedContext)
	}

	for _, account := range accounts {
		task.Log(fmt.Sprintf("processing account: %s", account.Name), "info", sharedContext)

		// get domains for account
		var domains []models.Domain
		sharedContext.PersistentDB.Where("account_id = ?", account.Id).Find(&domains)

		// get apps for account
		var apps []models.App
		sharedContext.PersistentDB.Where("account_id = ?", account.Id).Find(&apps)

		// get port mapping for account
		port_mapping := make(map[string]string, 10)

		for _, app := range apps {
			var image models.Image
			for _, i := range images {
				if i.Id == app.Image_id {
					image = i
					break
				}
			}

			if image.Id <= 0 {
				task.Log(fmt.Sprintf("image doesn't exist: %d (app: %d)", app.Image_id, app.Id), "error", sharedContext)
				continue
			}

			for _, port := range image.Ports {
				ip := "255.255.255.255"

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
					c, err := sharedContext.DockerClient.InspectContainer(con.ID)

					if err == nil {
						address := c.NetworkSettings.IPAddress

						if address != "" {
							ip = address
						}
					}
				}

				port_mapping[fmt.Sprintf("%s_%d_ip", app.Name, port.Port)] = ip
			}
		}

		vars, err := json.Marshal(port_mapping)
		if err == nil {
			task.Log(fmt.Sprintf("port mapping: %s", vars), "info", sharedContext)
		}

		for _, domain := range domains {
			task.Log(fmt.Sprintf("processing domain: %s", domain.Name), "info", sharedContext)

			nginx := ""
			apache := domain.Apache_config

			if domain.Redirect_url != "" {
				var result bytes.Buffer
				t, err := template.New("").Parse(`rewrite ^(.*) {{.Redirect_url}} permanent;
`)
				if err == nil {
					t.Execute(&result, domain)
					nginx = result.String()
				}
			} else {
				// parse nginx config
				nginx = domain.Nginx_config

				for k, v := range port_mapping {
					nginx = strings.Replace(nginx, fmt.Sprintf("#%s#", k), v, -1)
				}

				if domain.Apache_enabled {
					for k, v := range port_mapping {
						apache = strings.Replace(apache, fmt.Sprintf("#%s#", k), v, -1)
					}

					nginx = strings.Replace(nginx, "#apache#", "127.0.0.1:8080", -1)
				}
			}

			config := struct {
				Domain        models.Domain
				Nginx_Config  string
				Apache_Config string
			}{
				Domain:        domain,
				Nginx_Config:  nginx,
				Apache_Config: "",
			}

			t, err := template.New("").Parse(`server {
        {{ if .Domain.Ssl_enabled }}
        listen      443 spdy;
        {{ else }}
        listen 80;
        {{ end }}
        server_name {{.Domain.Name}};

        {{.Nginx_Config}}

        error_log /var/log/nginx/{{.Domain.Name}}-error.log warn;
        access_log /var/log/nginx/{{.Domain.Name}}-access.log combined;
}
`)

			if err == nil {
				var result bytes.Buffer

				t.Execute(&result, config)
				nginx = result.String()

				task.Log(fmt.Sprintf("nginx config: %s", nginx), "info", sharedContext)

				filename := fmt.Sprintf("%s_%s.conf", account.Name, domain.Name)
				if domain.Ssl_enabled {
					filename = fmt.Sprintf("%s_%s_ssl.conf", account.Name, domain.Name)
				}

				filepath := path.Join("/etc/nginx/manager", filename)

				err := ioutil.WriteFile(filepath, []byte(nginx), 0644)
				if err != nil {
					task.Log(fmt.Sprintf("error while writing nginx config: %s", err), "error", sharedContext)
					success = false
				}

				// test nginx config
				_, _, err = ExecuteShellCommandForTask("sudo", []string{"nginx", "-t"}, task, sharedContext)
				if err != nil {
					task.Log(fmt.Sprintf("nginx config failed (domain: %s): %s", domain.Name, err), "error", sharedContext)
					success = false

					// test failed, remove config file
					err := os.Remove(filepath)
					if err != nil {
						task.Log(fmt.Sprintf("error while removing nginx file (%s): %s", filepath, err), "error", sharedContext)
						success = false
					}
				}
			} else {
				task.Log(fmt.Sprintf("error while parsing template: %s", err), "error", sharedContext)
				success = false
			}

			if domain.Apache_enabled {
				// handle php tag
				r, _ := regexp.Compile("<php>(.*)</php>")

				name := strings.Replace(domain.Name, ".", "-", -1)

				config.Apache_Config = string(r.ReplaceAllFunc([]byte(apache), func(s []byte) []byte {
					port := strings.Replace(strings.Replace(string(s), "<php>", "", -1), "</php>", "", -1)

					return []byte(fmt.Sprintf(`AddType application/x-httpd-fastphp5 .php
                Action application/x-httpd-fastphp5 /php5-fcgi
                Alias /php5-fcgi /usr/lib/cgi-bin/php5-fcgi-%s
                FastCgiExternalServer /usr/lib/cgi-bin/php5-fcgi-%s -host %s:9000 -pass-header Authorization`,
						name, name, port))
				}))

				// handle apache config
				t, err := template.New("").Parse(`<VirtualHost *:8080>
        ServerName {{.Domain.Name}}

        {{.Apache_Config}}

        ErrorLog /var/log/apache2/{{.Domain.Name}}-error.log
        CustomLog /var/log/apache2/{{.Domain.Name}}-access.log combined
</VirtualHost>
`)

				if err == nil {
					var result bytes.Buffer

					// execute template
					t.Execute(&result, config)
					apache = result.String()

					task.Log(fmt.Sprintf("apache config: %s", apache), "info", sharedContext)

					filename := fmt.Sprintf("%s_%s.conf", account.Name, domain.Name)

					filepath := path.Join("/etc/apache2/manager", filename)

					// save config
					err := ioutil.WriteFile(filepath, []byte(apache), 0644)
					if err != nil {
						task.Log(fmt.Sprintf("error while writing apache config: %s", err), "error", sharedContext)
						success = false
					}

					// test apache config
					_, _, err = ExecuteShellCommandForTask("sudo", []string{"apachectl", "configtest"}, task, sharedContext)
					if err != nil {
						task.Log(fmt.Sprintf("apache config failed (domain: %s): %s", domain.Name, err), "error", sharedContext)
						success = false

						// test failed, remove config file
						err := os.Remove(filepath)
						if err != nil {
							task.Log(fmt.Sprintf("error while removing apache file (%s): %s", filepath, err), "error", sharedContext)
							success = false
						}
					}

				} else {
					task.Log(fmt.Sprintf("error while parsing template: %s", err), "error", sharedContext)
					success = false
				}
			}
		}

		_, _, err = ExecuteShellCommandForTask("sudo", []string{"service", "nginx", "reload"}, task, sharedContext)
		if err != nil {
			task.Log(fmt.Sprintf("nginx reload failed: %s", err), "error", sharedContext)
			success = false
		}

		_, _, err = ExecuteShellCommandForTask("sudo", []string{"service", "apache2", "reload"}, task, sharedContext)
		if err != nil {
			task.Log(fmt.Sprintf("nginx reload failed: %s", err), "error", sharedContext)
			success = false
		}
	}

	return success
}
