package main

import (
	"fmt"

	docs "github.com/dominiclet/golang-base/docs"
	initserver "github.com/dominiclet/golang-base/init_server"
	"github.com/dominiclet/golang-base/init_server/env"
	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	PORT = "8080"
)

// @title Golang base server
// @version 1.0

func main() {
	if !env.IsDevDirect() {
		gin.SetMode(gin.ReleaseMode)
	}

	logger := logger.InitLogger()

	router := initserver.InitDeps()

	r := gin.Default()

	router.RegisterRoutes(r)

	// Register path for swagger
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	logger.WithField("port", PORT).Info("Starting server")
	r.Run(fmt.Sprintf(":%s", PORT))
}
