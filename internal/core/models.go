package core

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type Namespace struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"`
}

type Project struct {
	ID                int        `json:"id"`
	Description       string     `json:"description"`
	SSHURLToRepo      string     `json:"ssh_url_to_repo"`
	WebURL            string     `json:"web_url"`
	Name              string     `json:"name"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreatedAt         *time.Time `json:"created_at"`
	LastActivityAt    *time.Time `json:"last_activity_at"`
	Namespace         *Namespace `json:"namespace"`
}

func mapNamespace(in *gitlab.ProjectNamespace) *Namespace {
	return &Namespace{
		ID:       in.ID,
		Name:     in.Name,
		Path:     in.Path,
		Kind:     in.Kind,
		FullPath: in.FullPath,
	}
}

func mapProject(in *gitlab.Project) *Project {
	return &Project{
		ID:                in.ID,
		Description:       in.Description,
		SSHURLToRepo:      in.SSHURLToRepo,
		WebURL:            in.WebURL,
		Name:              in.Name,
		Path:              in.Path,
		PathWithNamespace: in.PathWithNamespace,
		CreatedAt:         in.CreatedAt,
		LastActivityAt:    in.LastActivityAt,
		Namespace:         mapNamespace(in.Namespace),
	}
}
