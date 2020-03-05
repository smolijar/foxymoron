package core

import (
	"regexp"
	"strings"

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
	Count       int
	Types       map[string]int
	Issues      *Occurences
	Gitmoji     *Occurences
	IssuePrefix *Occurences
}

func isMergeCommit(commit *gitlab.Commit) bool {
	return len(commit.ParentIDs) > 1
}

func isSuggestion(commit *gitlab.Commit) bool {
	matcher := regexp.MustCompile(`^Apply suggestion to .*`)
	return matcher.MatchString(commit.Title)
}

func isRevert(commit *gitlab.Commit) bool {
	matcher := regexp.MustCompile(`^Revert ".*`)
	return matcher.MatchString(commit.Title)
}

type Occurences struct {
	Count      int
	Occurences map[string]int
}

func countOccurences(commit *gitlab.Commit, pattern *regexp.Regexp, result *Occurences) {
	matches := pattern.FindStringSubmatch(commit.Message)
	if matches != nil && len(matches) > 0 {
		opener := matches[1]
		result.Occurences[opener] = result.Occurences[opener] + 1
		result.Count++
	}
}

func CommitsToStats(commits []*gitlab.Commit) (stats Stats) {
	gitmojis, _ := gitmoji.Fetch()
	stats.Count = len(commits)
	stats.Types = map[string]int{"merge": 0, "suggestion": 0, "revert": 0, "human": 0}
	stats.Issues = &Occurences{Occurences: map[string]int{}}
	stats.IssuePrefix = &Occurences{Occurences: map[string]int{}}
	stats.Gitmoji = &Occurences{Occurences: map[string]int{}}
	for _, c := range commits {
		if isMergeCommit(c) {
			stats.Types["merge"]++
		} else if isSuggestion(c) {
			stats.Types["suggestion"]++
		} else if isRevert(c) {
			stats.Types["revert"]++
		} else {
			stats.Types["human"]++
			countOccurences(c, regexp.MustCompile(`(?:^|\n)[^\n]*(#[0-9]+)`), stats.Issues)
			countOccurences(c, regexp.MustCompile(`(?:^|\n)([^\n]*)#[0-9]+`), stats.IssuePrefix)
			countOccurences(c, regexp.MustCompile("^("+strings.Join(gitmojis, "|")+").*"), stats.Gitmoji)
		}
	}
	return
}
