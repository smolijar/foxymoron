package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/grissius/foxymoron/api"
)

// @title Foxymoron REST API
// @version 1.0
// @description API Proxy to GitLab

// @license.name MIT

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKey
// @in header
// @name Authorization

// @securityDefinitions.apikey GitLabURL
// @in header
// @name X-Gitlab-Url
func createEngine() *gin.Engine {
	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(authMdw)

	r.GET("/projects", getProjectsController)
	r.GET("/commits", getCommitsController)
	r.GET("/statistics", getStatisticsController)
	return r
}

func RunAt(port int) {
	log.Printf("Startig server on port %v", port)
	createEngine().Run(":" + strconv.Itoa(port))
}
