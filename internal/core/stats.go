package core

import (
	"log"
	"regexp"
	"sort"

	"github.com/grissius/foxymoron/pkg/gitmoji"
	"github.com/xanzy/go-gitlab"
)

type ProjectWithCommits struct {
	Project *gitlab.Project  `json:"project"`
	Commits []*gitlab.Commit `json:"commits"`
}

func groupByProject(projectsMap map[int]*gitlab.Project, commits []*gitlab.Commit) (res []*ProjectWithCommits) {
	projectsWithCommits := map[int][]*gitlab.Commit{}
	for _, c := range commits {
		if projectsWithCommits[c.ProjectID] == nil {
			projectsWithCommits[c.ProjectID] = []*gitlab.Commit{}
		}
		projectsWithCommits[c.ProjectID] = append(projectsWithCommits[c.ProjectID], c)
	}
	for pid, commits := range projectsWithCommits {
		res = append(res, &ProjectWithCommits{projectsMap[pid], commits})
	}
	return
}

type Stats struct {
	Count          int
	MergeCommits   int
	RefsPrefixes   map[string]int
	Issues         map[string]int
	Openers        map[string]int
	WithReferences int
	WithGitmoji    int
}

func CommitsToStats(commits []*gitlab.Commit) (stats Stats) {
	gitmojis, _ := gitmoji.Fetch()
	stats.Count = len(commits)
	stats.RefsPrefixes = map[string]int{}
	stats.Issues = map[string]int{}
	stats.Openers = map[string]int{}
	refernceMatcher := regexp.MustCompile(`^|\n(.+)(#[0-9]+)`)
	openerMatcher := regexp.MustCompile(`^(\S+)\s`)
	for _, c := range commits {
		stats.Count++
		if len(c.ParentIDs) > 1 {
			stats.MergeCommits++
		}
		refs := refernceMatcher.FindAllStringSubmatch(c.Message, -1)
		if refs != nil && len(refs) > 1 && len(refs[1]) > 2 {
			prefix := refs[1][1]
			stats.RefsPrefixes[prefix]++
			stats.WithReferences++
			issue := refs[1][2]
			stats.Issues[issue]++
			log.Println(prefix, issue)
		}

		openers := openerMatcher.FindAllStringSubmatch(c.Message, -1)
		if openers != nil && len(openers) > 0 && len(openers[0]) > 1 {
			opener := openers[0][1]
			stats.Openers[opener]++
			i := sort.SearchStrings(gitmojis, opener)
			if i < len(gitmojis) && gitmojis[i] == opener {
				stats.WithGitmoji++
			}
		}

	}
	return
}
