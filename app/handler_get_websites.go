package app

import (
	"github.com/boreq/plum/core"
)

type GetWebsites struct {
}

type GetWebsitesHandler struct {
	repositories *core.Repositories
}

func NewGetWebsitesHandler(repositories *core.Repositories) *GetWebsitesHandler {
	return &GetWebsitesHandler{
		repositories: repositories,
	}
}

func (h *GetWebsitesHandler) Execute(query GetWebsites) ([]string, error) {
	return h.repositories.Names(), nil
}
