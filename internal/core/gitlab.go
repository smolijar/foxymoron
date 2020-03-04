package core

import (
	"log"
	"regexp"
	"time"

	"github.com/xanzy/go-gitlab"
)

type User struct {
	GitLabURL   string
	ApiKey      string
	Client      *gitlab.Client
	ProjectsMap map[int]*Project
}

type FetchCommitsOptions struct {
	From         *time.Time
	To           *time.Time
	WithStats    bool
	MessageRegex *regexp.Regexp
}

func FetchCommits(user *User, opts *FetchCommitsOptions) []*gitlab.Commit {
	opt := &gitlab.ListCommitsOptions{
		Since:     opts.From,
		Until:     opts.To,
		All:       gitlab.Bool(opts.WithStats),
		WithStats: gitlab.Bool(true),
	}
	requests := 0
	requestsFinished := 0
	bound := make(map[int]bool)
	projects := user.ProjectsMap
	commitsChan := make(chan []*gitlab.Commit)
	for _, p := range projects {
		if p.LastActivityAt.Before(*opts.From) || p.CreatedAt.After(*opts.To) {
			continue
		}
		proj := p
		requests++
		go func() {
			commits, _, _ := user.Client.Commits.ListCommits(proj.ID, opt)
			for _, c := range commits {
				c.ProjectID = proj.ID
				bound[c.ProjectID] = true
			}
			commitsChan <- commits
			requestsFinished++
			log.Printf("Finished commit list request %d/%d", requestsFinished, requests)
		}()
	}
	commits := []*gitlab.Commit{}
	retrievedCommitsN := 0
	for i := 0; i < requests; i++ {
		for _, c := range <-commitsChan {
			retrievedCommitsN++
			if opts.MessageRegex == nil || opts.MessageRegex.MatchString(c.Message) {
				commits = append(commits, c)
			}
		}
	}
	log.Printf("Returning %v commits - Filtered from %v retrieved commits from %v projects for range <%v, %v>", retrievedCommitsN, len(commits), len(projects), opts.From, opts.To)
	log.Printf("GitLab request stats: %d (all projects) > %d (requests) > %d (lower bound for requests)", len(projects), requests, len(bound))
	return commits
}

func FetchProjects(user *User) (res []*Project) {
	for _, project := range user.ProjectsMap {
		res = append(res, project)
	}
	return
}

func FetchProjectsMap(user *User) map[int]*Project {
	maxPage := 0
	projectsMap := make(map[int]*Project)
	projectsChannel := make(chan []*gitlab.Project)
	maxPageChan := make(chan int)
	getProjectPage := func(currentPage int) {
		log.Printf("Making project list request %v / %v", currentPage, maxPage)
		ps, res, _ := user.Client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Simple: gitlab.Bool(true),
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    currentPage,
			},
		})
		select {
		case maxPageChan <- res.TotalPages:
		default:
		}
		projectsChannel <- ps
	}
	for page := 1; page == 1 || page <= maxPage; page++ {
		go getProjectPage(page)
		if page == 1 {
			maxPage = <-maxPageChan
		}
	}
	for page := 0; page < maxPage; page++ {
		for _, p := range <-projectsChannel {
			projectsMap[p.ID] = mapProject(p)
		}
	}
	log.Printf("Fetched %v projects from GitLab", len(projectsMap))
	return projectsMap
}
