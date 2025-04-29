package domain

import (
	"errors"
	"time"

	"github.com/bilte-co/bilte/internal/toolshed"
)

const FileChangeULIDPrefix = "flc"

type FileChange struct {
	ID             int
	ULID           string
	RepositoryID   int
	CommitID       int
	FilePath       string
	Language       string
	LinesAdded     int
	LinesRemoved   int
	LinesChanged   int
	CommitHash     string
	VendorFiles    bool
	GeneratedFiles bool
	CommitAt       *time.Time
	SyncedAt       *time.Time
}

func NewFileChange(repositoryID, commitID int, filePath, language string, linesAdded, linesRemoved, linesChanged int, commitHash string, commitAt time.Time) (*FileChange, error) {
	ulid, err := toolshed.CreateULID(FileChangeULIDPrefix, commitAt)
	if err != nil {
		return nil, errors.New("errors creating file change uild: " + err.Error())
	}

	return &FileChange{
		ULID:           ulid,
		RepositoryID:   repositoryID,
		CommitID:       commitID,
		FilePath:       filePath,
		Language:       language,
		LinesAdded:     linesAdded,
		LinesRemoved:   linesRemoved,
		LinesChanged:   linesChanged,
		CommitHash:     commitHash,
		VendorFiles:    false,
		GeneratedFiles: false,
		CommitAt:       &commitAt,
	}, nil
}
