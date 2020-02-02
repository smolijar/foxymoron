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

var projectsMap map[int]*gitlab.Project = nil

func fetchProjects() (res []*gitlab.Project) {
	for _, project := range fetchProjectsMap() {
		res = append(res, project)
	}
	return
}

func fetchProjectsMap() map[int]*gitlab.Project {
	// TODO: add mutex to prevent duplicate fetchers
	if projectsMap != nil {
		log.Printf("Got %v projects from cache", len(projectsMap))
		return projectsMap
	}
	newProjects := map[int]*gitlab.Project{}
	maxPage := 1
	projectsChannel := make(chan []*gitlab.Project)
	collectResults := func() {
		for _, p := range <-projectsChannel {
			newProjects[p.ID] = p
		}
	}
	for p := 1; true; p++ {
		go func(page int) {
			log.Printf("Making project list request %v/%v\n", page, maxPage)
			ps, res, _ := getClient().Projects.ListProjects(&gitlab.ListProjectsOptions{
				Archived: gitlab.Bool(false),
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
					Page:    p,
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
	projectsMap = newProjects
	log.Printf("Fetched %v projects from GitLab", len(projectsMap))
	return projectsMap
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
	projects := fetchProjectsMap()
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
	projectsWithCommits := map[int][]*gitlab.Commit{}
	for _, c := range commits {
		if projectsWithCommits[c.ProjectID] == nil {
			projectsWithCommits[c.ProjectID] = []*gitlab.Commit{}
		}
		projectsWithCommits[c.ProjectID] = append(projectsWithCommits[c.ProjectID], c)
	}
	projectsMap := fetchProjectsMap()
	for pid, commits := range projectsWithCommits {
		res = append(res, &ProjectWithCommits{projectsMap[pid], commits})
	}
	return
}
