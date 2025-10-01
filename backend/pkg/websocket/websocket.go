package websocket

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,

  CheckOrigin: func(r *http.Request) bool { return true },
}

// Upgrade upgrades the HTTP request to a Gorilla WebSocket connection.
// If the handshake fails, it logs the error and returns nil and the error.
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
			log.Println(err)
			return nil, err
	}
	return conn, nil
}