package api

import (
	"github.com/gin-gonic/gin"
	"github.com/grissius/foxymoron/internal/core"
	"github.com/xanzy/go-gitlab"
)

func getClient(c *gin.Context) *gitlab.Client {
	return c.MustGet("client").(*gitlab.Client)
}

func authMdw(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	gitlabUrl := c.GetHeader("X-Gitlab-Url")
	if authorization == "" || gitlabUrl == "" {
		c.AbortWithStatus(401)
	}
	c.Set("client", core.CreateClient(&authorization, &gitlabUrl))
}
