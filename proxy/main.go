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
var dockerClient *docker.Client

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

    reader, err := mydocker.ContainerLogs(c.Params.ByName("id"), &options)
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
            c.Set("uid", userId)
            log.Printf("OK")
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
            c.Set("account", account)
        }else{
            c.Fail(404, errors.New("Not Found"))
        }

        c.Next()
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

    //go-dockerclient
    endpoint := "unix:///var/run/docker.sock"
    dockerClient, _ = docker.NewClient(endpoint)

    //Listen to events
    listener := make(chan *docker.APIEvents)
    go dockerEvents(listener)
    dockerClient.AddEventListener(listener)

    sharedContext = &shared.SharedContext{}
    sharedContext.OpenDB()
    sharedContext.PersistentDB.AutoMigrate(&models.CronJob{})
    sharedContext.LogDB.AutoMigrate(&models.CronJobLog{})

    //HTTP
    r := gin.Default()

    r.LoadHTMLGlob("templates/*")

    authorized := r.Group("/")

    authorized.Use(Authentication())
    {
        authorized.GET("/ws/", func(c *gin.Context) {
            wsHandler(c.Writer, c.Request)
        })

        authorized.GET("/api/v1/containers/:id/logs", containerLogsHandler)

        authorized.GET("/api/v1/account/:account/shell", func(c *gin.Context) {
            shell.WebSocketShell(c)
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
            c.JSON(200, images)
        })

        //accounts
        accounts := &controllers.AccountsAPI{
            Context: sharedContext,
        }

        authorized.GET("/api/v1/accounts", accounts.ListAccounts)

        requiresAccount := authorized.Group("/api/v1/accounts/:name")

        requiresAccount.Use(RequireAccount())
        {
            requiresAccount.GET("", accounts.GetAccountByName)
            requiresAccount.GET("/apps", accounts.GetApps)

            //cronjobs
            cronJobs := &controllers.CronJobsAPI{
                Context: sharedContext,
            }

            requiresAccount.GET("/cronjobs", cronJobs.ListCronjobs)
            requiresAccount.GET("/cronjobs/:id", cronJobs.GetCronjob)
            requiresAccount.PUT("/cronjobs/:id", cronJobs.EditCronjob)
            requiresAccount.POST("/cronjobs", cronJobs.AddCronjob)
        }

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

    r.NoRoute(proxyRequest)

    //certFile := "../ssl.crt"
    //keyFile := "../ssl.key"

    /*if _, err := os.Stat(certFile); err == nil {
        log.Fatal(s.ListenAndServeTLS(certFile, keyFile))
    }else{*/
        log.Fatal(s.ListenAndServe())
    //}
}
