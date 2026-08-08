package app

import "github.com/boreq/plum/domain"

type GetWebsites struct {
}

type GetWebsitesHandler struct {
	repositories *domain.Repositories
}

func NewGetWebsitesHandler(repositories *domain.Repositories) *GetWebsitesHandler {
	return &GetWebsitesHandler{
		repositories: repositories,
	}
}

func (h *GetWebsitesHandler) Execute(query GetWebsites) ([]string, error) {
	return h.repositories.Names(), nil
}
