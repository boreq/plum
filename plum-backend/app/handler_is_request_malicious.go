package app

import (
	"time"

	"github.com/boreq/plum/plum-backend/domain"
	"github.com/boreq/plum/plum-backend/domain/request"
)

type IsRequestMalicious struct {
	RemoteAddress request.RemoteAddress
	Uri           request.Uri
	UserAgent     request.UserAgent
}

type IsRequestMaliciousHandler struct {
	maliciousAddresses MaliciousAddresses
	classifier         *domain.TrafficClassifier
	whitelist          domain.Whitelist
}

func NewIsRequestMaliciousHandler(maliciousAddresses MaliciousAddresses, classifier *domain.TrafficClassifier, whitelist domain.Whitelist) *IsRequestMaliciousHandler {
	return &IsRequestMaliciousHandler{
		maliciousAddresses: maliciousAddresses,
		classifier:         classifier,
		whitelist:          whitelist,
	}
}

func (h *IsRequestMaliciousHandler) Execute(query IsRequestMalicious) (bool, error) {
	now := time.Now()

	if h.whitelist.Contains(query.RemoteAddress) {
		return false, nil
	}

	if h.classifier.Classify(query.Uri, query.UserAgent, now) == domain.CategoryMalicious {
		return true, nil
	}

	return h.maliciousAddresses.IsIpMalicious(now, query.RemoteAddress), nil
}
