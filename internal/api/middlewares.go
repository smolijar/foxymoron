package api

import (
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
	"github.com/grissius/foxymoron/internal/core"
)

var getCache = func() func() (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return func() (*ristretto.Cache, error) {
		if err != nil {
			return nil, err
		}
		return cache, nil
	}
}()

func getUser(c *gin.Context) *core.User {
	return c.MustGet("user").(*core.User)
}

func authMdw(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	gitlabUrl := c.GetHeader("X-Gitlab-Url")
	if authorization == "" || gitlabUrl == "" {
		c.AbortWithStatus(401)
	}
	client := core.CreateClient(&authorization, &gitlabUrl)
	user := core.User{gitlabUrl, authorization, client, nil}

	cache, cacheErr := getCache()
	var projectsMap map[int]*core.Project
	if cacheErr == nil {
		cached, found := cache.Get(authorization + gitlabUrl)
		if found {
			log.Printf("Cache hit")
			projectsMap = cached.(map[int]*core.Project)
		}
	}
	if projectsMap == nil {
		projectsMap = core.FetchProjectsMap(&user)
	}
	if cacheErr == nil {
		cache.SetWithTTL(authorization+gitlabUrl, projectsMap, 1, time.Duration(10)*time.Minute)
	}
	user.ProjectsMap = projectsMap

	c.Set("user", &user)
}
