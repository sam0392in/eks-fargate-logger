package main

import (
	"eks-fargate-logger/controllers"
	"eks-fargate-logger/helper"
	"eks-fargate-logger/middleware"
	"github.com/gin-gonic/gin"
)

var logger = helper.Logger()

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(middleware.Firewall)

	r.POST("/", controllers.ProcessLogs) // Firehose POST the log record to this endpoint
	//r.GET("/", tests.EsSend) //Enable this endpoint for test sending data to elasticsearch
	go controllers.FlushBuffer() // Start a goroutine to continuously flush the data buffer
	logger.Info("server ready for incoming requests")
	err := r.Run("0.0.0.0:3000")
	if err != nil {
		logger.Fatal("Error starting the server: " + err.Error())
	}
}
