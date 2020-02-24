package api

import (
	"github.com/gin-gonic/gin"
	"github.com/grissius/foxymoron/internal/core"
)

func getUser(c *gin.Context) *core.User {
	return c.MustGet("user").(*core.User)
}

func authMdw(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	gitlabUrl := c.GetHeader("X-Gitlab-Url")
	if authorization == "" || gitlabUrl == "" {
		c.AbortWithStatus(401)
	}
	c.Set("user", &core.User{gitlabUrl, authorization, core.CreateClient(&authorization, &gitlabUrl)})
}
