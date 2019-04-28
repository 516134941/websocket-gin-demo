package main

import (
	"flag"
	"log"
	"net/http"

	"test/websocket-gin-demo/message-push/client"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8999", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := client.NewHub()
	go hub.Run()
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		client.ServeWs(hub, c.Writer, c.Request)
	})
	router.Run(":8999")
}
