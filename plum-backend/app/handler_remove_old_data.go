package app

import "time"

type RemoveOldData struct {
	Now time.Time
}

type RemoveOldDataHandler struct {
	repositories       Repositories
	maliciousAddresses MaliciousAddresses
}

func NewRemoveOldDataHandler(repositories Repositories, maliciousAddresses MaliciousAddresses) *RemoveOldDataHandler {
	return &RemoveOldDataHandler{
		repositories:       repositories,
		maliciousAddresses: maliciousAddresses,
	}
}

func (h *RemoveOldDataHandler) Execute(cmd RemoveOldData) error {
	h.repositories.RemoveOldData(cmd.Now)
	h.maliciousAddresses.RemoveOldData(cmd.Now)
	return nil
}
