package core

import (
	"log"
	"regexp"
	"time"

	"github.com/xanzy/go-gitlab"
)

type FetchCommitsOptions struct {
	From         *time.Time
	To           *time.Time
	WithStats    bool
	MessageRegex *regexp.Regexp
}

func FetchCommits(client *gitlab.Client, opts *FetchCommitsOptions) []*gitlab.Commit {
	opt := &gitlab.ListCommitsOptions{
		Since:     opts.From,
		Until:     opts.To,
		All:       gitlab.Bool(opts.WithStats),
		WithStats: gitlab.Bool(true),
	}
	projects := fetchProjectsMap(client)
	commitsChan := make(chan []*gitlab.Commit)
	for _, p := range projects {
		proj := p
		// COOL: create ad-hoc blocking-to-async functions
		go func() {
			commits, _, _ := client.Commits.ListCommits(proj.ID, opt)
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
			if opts.MessageRegex == nil || opts.MessageRegex.MatchString(c.Message) {
				commits = append(commits, c)
			}
		}
	}
	// COOL: you can use default logger from `log` and it outputs by default `2020/01/11 17:35:28 Retireved ...`
	// COOL: you can use %v for default formatting
	log.Printf("Returning %v commits - Filtered from %v retrieved commits from %v projects for range <%v, %v>", retrievedCommitsN, len(commits), len(projects), opts.From, opts.To)
	return commits
}

func FetchProjects(client *gitlab.Client) (res []*gitlab.Project) {
	for _, project := range fetchProjectsMap(client) {
		res = append(res, project)
	}
	return
}

func fetchProjectsMap(client *gitlab.Client) map[int]*gitlab.Project {
	projectsMap := map[int]*gitlab.Project{}
	maxPage := 1
	projectsChannel := make(chan []*gitlab.Project)
	collectResults := func() {
		for _, p := range <-projectsChannel {
			projectsMap[p.ID] = p
		}
	}
	for p := 1; true; p++ {
		go func(page int) {
			log.Printf("Making project list request %v/%v\n", page, maxPage)
			ps, res, _ := client.Projects.ListProjects(&gitlab.ListProjectsOptions{
				Archived: gitlab.Bool(false),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    page,
				},
			})
			maxPage = res.TotalPages
			projectsChannel <- ps
		}(p)
		if p == 1 {
			collectResults()
		}
		if p >= maxPage {
			break
		}

	}
	collectResults()
	log.Printf("Fetched %v projects from GitLab", len(projectsMap))
	return projectsMap
}
