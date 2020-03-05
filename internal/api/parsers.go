package api

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grissius/foxymoron/internal/core"
)

func parseIsoDate(date string) time.Time {
	t, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return time.Now()
	}
	return t
}

func parseCommitOptions(c *gin.Context) *core.FetchCommitsOptions {
	from := parseIsoDate(c.Query("from"))
	to := parseIsoDate(c.Query("to"))
	message, _ := regexp.Compile(c.Query("message"))
	return &core.FetchCommitsOptions{From: &from, To: &to, WithStats: true, MessageRegex: message}
}
