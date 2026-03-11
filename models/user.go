package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type User struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

func NewUser(conn *websocket.Conn) *User {
	return &User{
		ID:   generateID(),
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (u *User) ReadPump() {
	defer func() {
		hub := GetHub()
		hub.Unregister(u)
		u.Conn.Close()
	}()

	u.Conn.SetReadDeadline(time.Now().Add(pongWait))
	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := u.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msg.UserID = u.ID
		msg.Timestamp = time.Now()

		hub := GetHub()
		hub.Broadcast(&msg)
	}
}

func (u *User) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-u.Send:
			u.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(u.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-u.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			u.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
