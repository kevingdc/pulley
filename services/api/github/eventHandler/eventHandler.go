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
	case event.EventInstallation:
		return NewInstallationEventHandler(e)
	case event.EventPullRequest:
		return pr.New(e)
	default:
		return nil
	}
}
