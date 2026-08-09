package app

import "time"

type RemoveOldData struct {
	Now time.Time
}

type RemoveOldDataHandler struct {
	repositories Repositories
}

func NewRemoveOldDataHandler(repositories Repositories) *RemoveOldDataHandler {
	return &RemoveOldDataHandler{
		repositories: repositories,
	}
}

func (h *RemoveOldDataHandler) Execute(cmd RemoveOldData) error {
	h.repositories.RemoveOldData(cmd.Now)
	return nil
}
