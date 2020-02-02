package main

import (
	"fmt"
	"regexp"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/grissius/foxymoron/docs"
	"github.com/xanzy/go-gitlab"

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

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} gitlab.Commit "ok"
// @Header 200 {string} Token "qwerty"
// @Router /accounts/{id} [get]
func getCommitsController(c *gin.Context) {
	from := parseIsoDate(c.Query("from"))
	to := parseIsoDate(c.Query("to"))
	message, _ := regexp.Compile(c.Query("message"))
	commits := fetchCommits(&FetchCommitsOptions{from: &from, to: &to, withStats: true, messageRegex: message})
	c.JSON(200, commits)
}

func createRouter() *gin.Engine {
	fmt.Printf("%T", gitlab.Commit{})

	r := gin.Default()

	r.GET("/projects", getProjectsController)
	r.GET("/commits", getCommitsController)
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
