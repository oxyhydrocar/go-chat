package models

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	users      map[*User]bool
	broadcast  chan *Message
	register   chan *User
	unregister chan *User
	mu         sync.RWMutex
}

var (
	instance *Hub
	once     sync.Once
)

func GetHub() *Hub {
	once.Do(func() {
		instance = &Hub{
			broadcast:  make(chan *Message),
			register:   make(chan *User),
			unregister: make(chan *User),
			users:      make(map[*User]bool),
		}
		go instance.run()
	})
	return instance
}

func (h *Hub) run() {
	for {
		select {
		case user := <-h.register:
			h.mu.Lock()
			h.users[user] = true
			h.mu.Unlock()

		case user := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.users[user]; ok {
				delete(h.users, user)
				close(user.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			messageBytes, _ := json.Marshal(message)
			for user := range h.users {
				select {
				case user.Send <- messageBytes:
				default:
					close(user.Send)
					delete(h.users, user)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(user *User) {
	h.register <- user
}

func (h *Hub) Unregister(user *User) {
	h.unregister <- user
}

func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

func (h *Hub) GetActiveUsers() []*User {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]*User, 0, len(h.users))
	for user := range h.users {
		users = append(users, user)
	}
	return users
}
