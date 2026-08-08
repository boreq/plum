package app

import (
	"time"

	"github.com/boreq/plum/domain"
)

type RemoveOldData struct {
	Now time.Time
}

type RemoveOldDataHandler struct {
	repositories *domain.Repositories
}

func NewRemoveOldDataHandler(repositories *domain.Repositories) *RemoveOldDataHandler {
	return &RemoveOldDataHandler{
		repositories: repositories,
	}
}

func (h *RemoveOldDataHandler) Execute(cmd RemoveOldData) error {
	h.repositories.RemoveOldData(cmd.Now)
	return nil
}
