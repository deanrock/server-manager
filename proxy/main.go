package main
 
import (
    "log"
    "net/http"
    "net/http/httputil"
    "github.com/gorilla/websocket"
    "sync"
    "time"
    "bufio"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/samalba/dockerclient"
)


// Docker
func dockerEventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
    log.Printf("Received event: %#v\n", *event)

    broadcast(fmt.Sprintf("Received event: %#v\n", *event))
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

var docker *dockerclient.DockerClient

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

    reader, err := docker.ContainerLogs(c.Params.ByName("id"), &options)
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
            log.Printf("OK")
        }else{
            c.Fail(401, errors.New("Unauthorized"))
        }

        c.Next()
    }
}

 
func main() {
    //Docker
    // Init the client
    docker, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

    // Listen to events
    docker.StartMonitorEvents(dockerEventCallback, nil)


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

    r.NoRoute(proxyRequest)
    log.Fatal(s.ListenAndServe())
}
