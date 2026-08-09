package app

import (
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
)

type AddRequest struct {
	Website domain.WebsiteName
	Entry   *parser.Entry
}

type AddRequestHandler struct {
	repositories Repositories
	classifier   *domain.TrafficClassifier
}

func NewAddRequestHandler(repositories Repositories, classifier *domain.TrafficClassifier) *AddRequestHandler {
	return &AddRequestHandler{
		repositories: repositories,
		classifier:   classifier,
	}
}

func (h *AddRequestHandler) Execute(cmd AddRequest) error {
	repository, ok := h.repositories.Get(cmd.Website)
	if !ok {
		return ErrWebsiteNotFound
	}

	category := h.classifier.Classify(cmd.Entry)

	return repository.Insert(cmd.Entry, category)
}
