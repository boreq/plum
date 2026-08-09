package app

import (
	"github.com/boreq/plum/domain"
	"github.com/boreq/plum/domain/parser"
)

type AddRequest struct {
	Website domain.WebsiteName
	Entry   *parser.Entry
}

type AddRequestHandler struct {
	repositories Repositories
}

func NewAddRequestHandler(repositories Repositories) *AddRequestHandler {
	return &AddRequestHandler{
		repositories: repositories,
	}
}

func (h *AddRequestHandler) Execute(cmd AddRequest) error {
	repository, ok := h.repositories.Get(cmd.Website)
	if !ok {
		return ErrWebsiteNotFound
	}

	return repository.Insert(cmd.Entry)
}
