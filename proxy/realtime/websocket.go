package realtime

import (
	//"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Connections struct {
	sync.RWMutex
	m map[*websocket.Conn]WSConnection
}

type WebsocketHandler struct {
	Connections Connections
}

type WSConnection struct {
	UserId      int
	UserIsStuff bool
}

func NewWebsocketHandler() WebsocketHandler {
	w := WebsocketHandler{}
	w.Connections = Connections{
		m: make(map[*websocket.Conn]WSConnection),
	}

	return w
}

// HTTP
func (w WebsocketHandler) Broadcast(message string) {
	w.Connections.Lock()
	for conn, _ := range w.Connections.m {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			delete(w.Connections.m, conn)
			conn.Close()
		}
	}
	w.Connections.Unlock()
}

func (w WebsocketHandler) SendToUser(msg []byte, userId int) {
	w.Connections.Lock()
	for conn, e := range w.Connections.m {
		if e.UserId == userId {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				delete(w.Connections.m, conn)
				conn.Close()
			}
		}
	}
}

func (w WebsocketHandler) WsHandler(c *gin.Context, userId int, userIsStuff bool, callback func(uid int)) {
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	wsconn := WSConnection{
		UserId:      userId,
		UserIsStuff: userIsStuff,
	}

	w.Connections.Lock()
	w.Connections.m[conn] = wsconn
	w.Connections.Unlock()

	callback(userId)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			w.Connections.Lock()
			delete(w.Connections.m, conn)
			w.Connections.Unlock()
			conn.Close()
			return
		}
	}
}
