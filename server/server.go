package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var messageChan = make(chan string)

//此处为服务端接收 需要推送的信息
func updateMsg() {
	for {
		var messageBuf string
		fmt.Scanln(&messageBuf)
		messageChan <- messageBuf
	}
}

//此处为websocket 执行的长连接 函数
func echo(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		//信息推送 WriteMessage
		//updateMsg()
		c.WriteMessage(websocket.TextMessage, []byte(<-messageChan))
	}
}

func main() {
	go updateMsg()
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		echo(c.Writer, c.Request)
	})
	router.Run(":8999")
}
