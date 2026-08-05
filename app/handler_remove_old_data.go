package app

import (
	"time"

	"github.com/boreq/plum/core"
)

type RemoveOldData struct {
	Now time.Time
}

type RemoveOldDataHandler struct {
	repositories *core.Repositories
}

func NewRemoveOldDataHandler(repositories *core.Repositories) *RemoveOldDataHandler {
	return &RemoveOldDataHandler{
		repositories: repositories,
	}
}

func (h *RemoveOldDataHandler) Execute(cmd RemoveOldData) error {
	h.repositories.RemoveOldData(cmd.Now)
	return nil
}
