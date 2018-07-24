/**
 * Created by IDEA.
 * User: godsoul
 */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/DeanThompson/ginpprof"
	"github.com/guobin8205/api_demo/utils/config"
	"net/http"
	"time"
	"github.com/guobin8205/api_demo/routers"
	_ "github.com/guobin8205/api_demo/models"
	"github.com/guobin8205/api_demo/middleware"
)

func main() {
	router := gin.New()
	router.Static("/data", "./data")
	//middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ParseRequest())
	if conf.RunMode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	}
	ginpprof.Wrap(router)
	routers.SetRouters(router)

	s := &http.Server{
		Addr:           conf.HttpPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
