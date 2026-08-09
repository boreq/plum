package app

import "github.com/boreq/plum/domain"

type GetWebsites struct {
}

type GetWebsitesHandler struct {
	repositories Repositories
}

func NewGetWebsitesHandler(repositories Repositories) *GetWebsitesHandler {
	return &GetWebsitesHandler{
		repositories: repositories,
	}
}

func (h *GetWebsitesHandler) Execute(query GetWebsites) ([]domain.WebsiteName, error) {
	return h.repositories.Names(), nil
}
