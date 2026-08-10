package app

import (
	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/request"
)

type AddRequest struct {
	Website domain.WebsiteName
	Request request.Request
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

	category := h.classifier.Classify(cmd.Request.Uri(), cmd.Request.UserAgent(), cmd.Request.Timestamp())

	h.maliciousAddresses.Insert(cmd.Request, category)

	return repository.Insert(cmd.Request, category)
}
