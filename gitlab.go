package main

import (
	"log"
	"regexp"
	"time"

	"github.com/xanzy/go-gitlab"
)

func getClient() *gitlab.Client {
	git := gitlab.NewClient(nil, config.token)
	git.SetBaseURL(config.url)
	return git
}

var projects []*gitlab.Project = nil

func fetchProjects() []*gitlab.Project {
	// TODO: add mutex to prevent duplicate fetchers
	if projects != nil {
		log.Println("Getting projects from cache")
		return projects
	}
	newProjects := []*gitlab.Project{}
	for p := 1; true; p++ {
		ps, res, _ := getClient().Projects.ListProjects(&gitlab.ListProjectsOptions{
			Archived: gitlab.Bool(false),
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    p,
			},
		})
		for _, p := range ps {
			newProjects = append(newProjects, p)
		}
		if p >= res.TotalPages {
			break
		}

	}
	projects = newProjects
	log.Printf("Fetched %v projects from GitLab", len(projects))
	return projects
}

type FetchCommitsOptions struct {
	from         *time.Time
	to           *time.Time
	withStats    bool
	messageRegex *regexp.Regexp
}

func fetchCommits(opts *FetchCommitsOptions) []*gitlab.Commit {
	opt := &gitlab.ListCommitsOptions{
		Since:     opts.from,
		Until:     opts.to,
		All:       gitlab.Bool(opts.withStats),
		WithStats: gitlab.Bool(true),
	}
	projects := fetchProjects()
	commitsChan := make(chan []*gitlab.Commit)
	for _, p := range projects {
		proj := p
		// COOL: create ad-hoc blocking-to-async functions
		go func() {
			commits, _, _ := getClient().Commits.ListCommits(proj.ID, opt)
			for _, c := range commits {
				c.ProjectID = proj.ID
			}
			commitsChan <- commits
		}()
	}
	commits := []*gitlab.Commit{}
	retrievedCommitsN := 0
	for i := 0; i < len(projects); i++ {
		retrievedCommitsN++
		// COOL: use `<-commitsChan` like an expression without assignment
		for _, c := range <-commitsChan {
			if opts.messageRegex.MatchString(c.Message) {
				commits = append(commits, c)
			}
		}
	}
	// COOL: you can use default logger from `log` and it outputs by default `2020/01/11 17:35:28 Retireved ...`
	// COOL: you can use %v for default formatting
	log.Printf("Returning %v commits - Filtered from %v retrieved commits from %v projects for range <%v, %v>", retrievedCommitsN, len(commits), len(projects), opts.from, opts.to)
	return commits
}

type ProjectWithCommits struct {
	Project *gitlab.Project  `json:"project"`
	Commits []*gitlab.Commit `json:"commits"`
}

func groupByProject(commits []*gitlab.Commit) (res []*ProjectWithCommits) {
	projects := map[int][]*gitlab.Commit{}
	for _, c := range commits {
		if projects[c.ProjectID] == nil {
			projects[c.ProjectID] = []*gitlab.Commit{}
		}
		projects[c.ProjectID] = append(projects[c.ProjectID], c)
	}
	for k, v := range projects {
		// TODO include whole project
		res = append(res, &ProjectWithCommits{&gitlab.Project{ID: k}, v})
	}
	return
}
