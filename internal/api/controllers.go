package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grissius/foxymoron/internal/core"
)

// List projects
// @Tags Projects
// @Summary List all available projects
// @Produce json
// @Success 200 {array} core.Project
// @Router /projects [get]
// @Security ApiKey
// @Security GitLabURL
func getProjectsController(c *gin.Context) {
	projects := core.FetchProjects(getUser(c))
	c.JSON(200, projects)
}

// List commits
// @Tags Commits
// @Summary List commit from all available projects within range
// @Produce json
// @Success 200 {array} gitlab.Commit
// @Router /commits [get]
// @Param from query string false "Include only commits newer than this, e.g. `2020-02-19T00:00:00.000Z`"
// @Param to query string false "Include only commits older than this, e.g. `2020-02-20T00:00:00.000Z`"
// @Param message query string false "Pass only commits matching this regex pattern, e.g. `foo|bar`"
// @Security ApiKey
// @Security GitLabURL
func getCommitsController(c *gin.Context) {
	commits := core.FetchCommits(getUser(c), parseCommitOptions(c))
	c.JSON(200, commits)
}

// Commit statistics
// @Tags Statistics
// @Summary Get statistics for commits within range
// @Produce json
// @Success 200 {array} core.Stats
// @Router /statistics [get]
// @Param from query string false "Include only commits newer than this, e.g. `2020-02-19T00:00:00.000Z`"
// @Param to query string false "Include only commits older than this, e.g. `2020-02-20T00:00:00.000Z`"
// @Param message query string false "Pass only commits matching this regex pattern, e.g. `foo|bar`"
// @Param mode query int false "Group by nothing (0), project (1), namespace (2)"
// @Security ApiKey
// @Security GitLabURL
func getStatisticsController(c *gin.Context) {
	user := getUser(c)
	mode := 0
	stringMode, ok := c.GetQuery("mode")
	if ok {
		mode, _ = strconv.Atoi(stringMode)
	}
	stats := core.CommitsToBuckets(core.FetchCommits(user, parseCommitOptions(c)), user.ProjectsMap, mode)

	c.JSON(200, stats)
}

func root(c *gin.Context) {
	c.JSON(200, struct {
		Now               time.Time
		WhatDoesTheFoxSay string
	}{time.Now(), "🦊"})
}
