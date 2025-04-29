package domain

import (
	"errors"
	"time"

	"github.com/bilte-co/bilte/internal/toolshed"
)

const RepoULIDPrefix = "repo"

type Repo struct {
	ID            int
	ULID          string
	RepoName      string
	RepoFullName  string
	RepoOwner     string
	DefaultBranch string
	Url           string
	SyncedAt      *time.Time
}

func NewRepo(repoName string, repoFullName string, repoOwner string, defaultBranch string, url string) (*Repo, error) {
	ulid, err := toolshed.CreateULID(RepoULIDPrefix, time.Now())
	if err != nil {
		return nil, errors.New("errors creating repo uild: " + err.Error())
	}

	return &Repo{
		ULID:          ulid,
		RepoName:      repoName,
		RepoFullName:  repoFullName,
		RepoOwner:     repoOwner,
		DefaultBranch: defaultBranch,
		Url:           url,
	}, nil
}
