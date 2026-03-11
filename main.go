package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go-chat/handlers"
	"go-chat/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/api/users", handlers.GetUsersHandler)
	router.HandleFunc("/api/messages", handlers.GetMessagesHandler)

	staticDir := filepath.Join("static", "assets")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(staticDir))))

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	user := models.NewUser(conn)
	hub := models.GetHub()
	hub.Register(user)

	go user.WritePump()
	user.ReadPump()
}

func normalizePath(path string) string {
	cleaned := filepath.Clean(path)
	normalized := strings.TrimSpace(cleaned)
	return normalized
}
