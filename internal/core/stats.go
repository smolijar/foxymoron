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

type Bucket struct {
	Project   *Project   `json:"project"`
	Namespace *Namespace `json:"namespace"`
	commits   []*gitlab.Commit
	Stats     Stats `json:"stats"`
}

// mode 0 - no bucketing (single master bucket)
// mode 1 - per project
// mode 2 - per namespace
func CommitsToBuckets(commits []*gitlab.Commit, projectsMap map[int]*Project, mode int) (res []*Bucket) {
	buckets := make(map[int]*Bucket)
	if mode == 0 {
		buckets[0] = &Bucket{commits: commits}
	} else if mode == 1 {
		for _, c := range commits {
			if buckets[c.ProjectID] == nil {
				buckets[c.ProjectID] = &Bucket{Project: projectsMap[c.ProjectID], commits: []*gitlab.Commit{}}
			}
			buckets[c.ProjectID].commits = append(buckets[c.ProjectID].commits, c)
		}
	} else if mode == 2 {
		for _, c := range commits {
			nsId := projectsMap[c.ProjectID].Namespace.ID
			if buckets[nsId] == nil {
				buckets[nsId] = &Bucket{Namespace: projectsMap[c.ProjectID].Namespace, commits: []*gitlab.Commit{}}
			}
			buckets[nsId].commits = append(buckets[nsId].commits, c)
		}
	}
	for _, b := range buckets {
		b.Stats = CommitsToStats(b.commits)
		res = append(res, b)
	}
	return
}
