package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func getProjectsController(c *gin.Context) {
	projects := fetchProjects()
	c.JSON(200, projects)
}

func parseIsoDate(date string) time.Time {
	t, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return time.Now()
	}
	return t
}

func getCommitsController(c *gin.Context) {
	from := parseIsoDate(c.Query("from"))
	to := parseIsoDate(c.Query("to"))
	message, _ := regexp.Compile(c.Query("message"))
	commits := fetchCommits(&FetchCommitsOptions{from: &from, to: &to, withStats: true, messageRegex: message})
	c.JSON(200, commits)
}

func getStatisticsController(c *gin.Context) {
	from := parseIsoDate(c.Query("from"))
	to := parseIsoDate(c.Query("to"))
	stats := commitsToStats(fetchCommits(&FetchCommitsOptions{from: &from, to: &to, withStats: true, messageRegex: nil}))

	c.JSON(200, stats)
}

func createRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/projects", getProjectsController)
	r.GET("/commits", getCommitsController)
	r.GET("/statistics", getStatisticsController)
	return r
}
