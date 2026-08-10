package app

import (
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/parser"
	"github.com/boreq/plum/plum-backend/domain/request"
)

type AddRequest struct {
	Website domain.WebsiteName
	Entry   *parser.Entry
}

type AddRequestHandler struct {
	repositories       Repositories
	maliciousAddresses MaliciousAddresses
	classifier         *domain.TrafficClassifier
}

func NewAddRequestHandler(repositories Repositories, maliciousAddresses MaliciousAddresses, classifier *domain.TrafficClassifier) *AddRequestHandler {
	return &AddRequestHandler{
		repositories:       repositories,
		maliciousAddresses: maliciousAddresses,
		classifier:         classifier,
	}
}

func (h *AddRequestHandler) Execute(cmd AddRequest) error {
	repository, ok := h.repositories.Get(cmd.Website)
	if !ok {
		return ErrWebsiteNotFound
	}

	category := h.classifier.Classify(
		request.NewUri(cmd.Entry.HttpRequestURI),
		request.NewUserAgent(cmd.Entry.UserAgent),
		cmd.Entry.Time,
	)

	h.maliciousAddresses.Insert(cmd.Entry, category)

	return repository.Insert(cmd.Entry, category)
}
