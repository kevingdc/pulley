package eventhandler

import (
	"github.com/kevingdc/pulley/services/api/github/event"
	"github.com/kevingdc/pulley/services/api/github/eventHandler/pr"
)

func Handle(e *event.Payload) (event.HandlerResponse, error) {
	eventHandler := resolve(e)

	if eventHandler == nil {
		return nil, nil
	}

	return eventHandler.Handle()
}

func resolve(e *event.Payload) event.Handler {
	switch e.Type {
	case event.TypeInstallation:
		return NewInstallationEventHandler(e)
	case event.TypePullRequest:
		return pr.New(e)
	case event.TypePullRequestReview:
		return NewPullRequestReviewEventHandler(e)
	default:
		return nil
	}
}
