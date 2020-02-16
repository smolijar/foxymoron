package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createEngine() *gin.Engine {
	r := gin.Default()

	r.GET("/projects", getProjectsController)
	r.GET("/commits", getCommitsController)
	r.GET("/statistics", getStatisticsController)
	return r
}

func RunAt(port int) {
	log.Printf("Startig server on port %v", port)
	createEngine().Run(":" + strconv.Itoa(port))
}
