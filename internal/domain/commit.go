package domain

import (
	"errors"
	"time"

	"github.com/bilte-co/bilte/internal/toolshed"
)

const CommitULIDPrefix = "cmt"

type Commit struct {
	ID           int
	ULID         string
	CommitHash   string
	RepositoryID int
	AuthorName   string
	AuthorEmail  string
	Message      string
	CommitAt     *time.Time
	SyncedAt     *time.Time
}

func NewCommit(commitHash string, repositoryID int, authorName string, authorEmail string, message string, commitAt time.Time) (*Commit, error) {
	ulid, err := toolshed.CreateULID(CommitULIDPrefix, commitAt)
	if err != nil {
		return nil, errors.New("errors creating commit uild: " + err.Error())
	}

	return &Commit{
		ULID:         ulid,
		CommitHash:   commitHash,
		RepositoryID: repositoryID,
		AuthorName:   authorName,
		AuthorEmail:  authorEmail,
		Message:      message,
		CommitAt:     &commitAt,
	}, nil
}
