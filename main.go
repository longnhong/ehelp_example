package main

import (
	"github.com/gin-gonic/gin"
	// 1. init first
	_ "ehelp/init"
	// 2. iniit 2nd
	"ehelp/api"
	"ehelp/cache"
	"ehelp/common"
	"ehelp/middleware"
	"ehelp/room"
	"ehelp/system"
	"net/http"
)

func main() {
	router := gin.New()
	//static
	router.StaticFS("/static", http.Dir("./upload"))
	router.StaticFS("/admin", http.Dir("./admin")).Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html")
		c.Next()
	})
	router.Use(middleware.AddHeader(), gin.Logger(), middleware.Recovery())

	router.StaticFS("/app", http.Dir("./app"))
	system.Launch()
	var timer, _ = common.NewDailyTimer("23:00", func() {
		cache.MoveOrderToOpen()
	})

	timer.Start()
	var timer2, _ = common.NewDailyTimer("01:00", func() {
		system.SetCacheOrderDay()
	})
	timer2.Start()

	//api
	rootAPI := router.Group("/api")
	api.InitApi(rootAPI)
	//	ws
	room.NewRoomServer(router.Group("/room"))
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
