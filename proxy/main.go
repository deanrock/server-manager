package main
 
import (
    "log"
    "net/http"
    "net/http/httputil"
    "github.com/gorilla/websocket"
    "sync"
    //"os"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "time"
    "bufio"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/samalba/dockerclient"
    "github.com/fsouza/go-dockerclient"
    //"database/sql"
    //"github.com/coopernurse/gorp"
    //"github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "./controllers"
    "./shared"
    "./models"
    "./shell"
    "strconv"
)


var sharedContext *shared.SharedContext

// Docker
func dockerEvents(listener chan *docker.APIEvents) {
    for {
        event := <- listener
        log.Printf("Received event: %#v\n", *event)

        broadcast(fmt.Sprintf("Received event: %#v\n", *event))
    }
}


// HTTP
func broadcast(message string) {
    connections.Lock()
    for conn, _ := range connections.m {
        if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
            delete(connections.m, conn)
            conn.Close()
        }
    }
    connections.Unlock()
}

var connections = struct {
    sync.RWMutex
    m map[*websocket.Conn]bool
}{m: make(map[*websocket.Conn]bool)}

var mydocker *dockerclient.DockerClient

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }
    log.Println("Succesfully upgraded connection")

    connections.Lock()
    connections.m[conn] = true
    connections.Unlock()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            connections.Lock()
            delete(connections.m, conn)
            connections.Unlock()
            conn.Close()
            return
        }
    }
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
        Follow: true,
        Stdout: true,
        Stderr: true,
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

    for _, app := range(a.Apps()) {
        if app.Id == int(id) {
            log.Println(app.Id)
            log.Println(app)

            containers, _ := sharedContext.DockerClient.ListContainers(docker.ListContainersOptions{})

            name := fmt.Sprintf("/%s", app.ContainerName(a.Name))
            for _, c := range(containers) {
                for _, n := range(c.Names) {
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
                c.Fail(500, errors.New("cannot handle"))
                c.Abort()
            }

            c.Set("uid", userId)
            c.Set("user", *u)
            
            var access []models.UserAccess
            sharedContext.PersistentDB.Where("user_id = ?", userId).Find(&access)
            c.Set("userAccess", access)
        }else{
            //c.Fail(401, errors.New("Unauthorized"))
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
            for _, ua := range(userAccess) {
                if ua.Account_id == account.Id {
                    found = true
                }
            }

            if found {
                c.Set("account", account)
            }else{
                c.Fail(401, errors.New("unauthorized access to account"))
            }
        }else{
            c.Fail(404, errors.New("Not Found"))
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
            for _, ua := range(userAccess) {
                if ua.Account_id == account.Id {
                    if accessName == "shell_access" && ua.ShellAccess == true {
                        found = true
                    }else if accessName == "app_access" && ua.AppAccess == true {
                        found = true
                    }else if accessName == "cronjob_access" && ua.CronjobAccess == true {
                        found = true
                    }else if accessName == "domain_access" && ua.DomainAccess == true {
                        found = true
                    }else if accessName == "ssh_access" && ua.SshAccess == true {
                        found = true
                    }else if accessName == "database_access" && ua.DatabaseAccess == true {
                        found = true
                    }
                }
            }

            if !found {
                c.Fail(401, errors.New("user doesnt have access to this resource"))
            }
        }else{
            c.Fail(404, errors.New("Not Found"))
        }

        c.Next()
    }
}

func RequireStaff() gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(models.User)

        if user.Is_staff {
           c.Next()
        }else{
            c.Fail(401, errors.New("Unauthorized"))
            c.Abort()
        }
    }
}

type Config struct {
    Server_name string `json:"server_name"`
    Mysql_connection_string string `json:"mysql_connection_string"`
}

type Profile struct {
    Server_name string `json:"server_name"`
    User models.User `json:"user"`
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
    sharedContext.OpenDB()
    sharedContext.PersistentDB.LogMode(true)
    sharedContext.PersistentDB.AutoMigrate(&models.CronJob{})
    sharedContext.PersistentDB.AutoMigrate(&models.UserAccess{})
    sharedContext.PersistentDB.Model(&models.UserAccess{}).AddUniqueIndex("idx_user_account", "user_id", "account_id")
    sharedContext.LogDB.AutoMigrate(&models.CronJobLog{})

    //go-dockerclient
    endpoint := "unix:///var/run/docker.sock"
    sharedContext.DockerClient, _ = docker.NewClient(endpoint)

    //Listen to events
    listener := make(chan *docker.APIEvents)
    go dockerEvents(listener)
    sharedContext.DockerClient.AddEventListener(listener)

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
            wsHandler(c.Writer, c.Request)
        })

        authorized.GET("/api/v1/profile", func (c *gin.Context) {
            uid := c.MustGet("uid").(*int)
            user, _ := models.FindUserById(sharedContext, *uid)
            p := Profile{
                Server_name:     config.Server_name,
                User:            *user,
            }

            c.JSON(200, p)
        })

        authorized.GET("/api/v1/shells", func(c *gin.Context) {
            s := shell.Shell{}
            s.GetDockerImages()
            c.JSON(200, s.ShellImages)
        })

        //images
        authorized.GET("/api/v1/images", func(c *gin.Context) {
            var images []models.Image
            sharedContext.PersistentDB.Find(&images)

            var ports []models.ImagePort
            sharedContext.PersistentDB.Find(&ports)

            for k, i := range(images) {
                for _, v := range(ports) {
                    if v.Image_id == i.Id {
                        images[k].Ports = append(images[k].Ports, v)
                    }
                }
            }

            c.JSON(200, images)
        })

        //accounts
        accounts := &controllers.AccountsAPI{
            Context: sharedContext,
        }

        authorized.GET("/api/v1/accounts", accounts.ListAccounts)
        authorized.GET("/api/v1/all-accounts", RequireStaff(), accounts.ListAllAccounts)

        requiresAccount := authorized.Group("/api/v1/accounts/:name")

        requiresAccount.Use(RequireAccount())
        {
            requiresAccount.GET("", accounts.GetAccountByName)

            //apps
            requiresAccount.GET("/apps", RequireUserAccess("shell_access"), accounts.GetApps)
            requiresAccount.GET("/apps/:id/logs", RequireUserAccess("shell_access"), containerLogsHandler)

            //shell
            requiresAccount.GET("/shell", RequireUserAccess("shell_access"), func(c *gin.Context) {
                shell.WebSocketShell(c, sharedContext)
            })

            //cronjobs
            cronJobs := &controllers.CronJobsAPI{
                Context: sharedContext,
            }

            requiresAccount.GET("/cronjobs", RequireUserAccess("cronjob_access"), cronJobs.ListCronjobs)
            requiresAccount.GET("/cronjobs/:id", RequireUserAccess("cronjob_access"), cronJobs.GetCronjob)
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
    }

    s := &http.Server{
        Addr:           ":4444",
        Handler:        r,
        ReadTimeout:    60 * time.Minute,
        WriteTimeout:   60 * time.Minute,
        MaxHeaderBytes: 1 << 20,
    }

    r.Static("/static/ace", "../manager/static/ace")
    r.Static("/static/bootstrap", "../manager/static/bootstrap")
    r.Static("/static/css", "../manager/static/css")
    r.Static("/static/js", "../manager/static/js")
    r.Static("/static/templates", "../manager/static/templates")
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
