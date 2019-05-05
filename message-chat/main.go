package main

import (
	"flag"
	"runtime"
	"test/websocket-gin-demo/message-chat/config"
	"test/websocket-gin-demo/message-chat/middlewares"
	"test/websocket-gin-demo/message-chat/server"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
)

var (
	tomlFile = flag.String("config", "docs/test.toml", "config file")
)

// init 初始化配置
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gin.SetMode(gin.DebugMode)
}
func main() {
	flag.Parse()
	// 解析配置文件
	tomlConfig, err := config.UnmarshalConfig(*tomlFile)
	if err != nil {
		log.Printf("UnmarshalConfig: err:%v\n", err)
		return
	}
	// 绑定路由，及公共的tomlConfig
	//router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery())
	//设置跨域
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"x-xq5-jwt", "Content-Type", "Origin", "Content-Length"},
		ExposeHeaders:    []string{"x-xq5-jwt"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middlewares.Config(tomlConfig))
	// router.Use(middlewares.Redis("test", tomlConfig))
	hub := server.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) { server.ServeWs(hub, c) })
	log.Println("run websocket at ", tomlConfig.GetListenAddr())
	router.Run(tomlConfig.GetListenAddr())
}
