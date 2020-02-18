package core

import "github.com/xanzy/go-gitlab"

func CreateClient(token *string, url *string) *gitlab.Client {
	git := gitlab.NewClient(nil, *token)
	git.SetBaseURL(*url)
	return git
}
