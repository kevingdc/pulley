package event

type EventHandlerResponse interface{}

type EventHandler interface {
	Handle() (EventHandlerResponse, error)
}

func Handle(e *EventPayload) (EventHandlerResponse, error) {
	eventHandler := resolve(e)

	if eventHandler == nil {
		return nil, nil
	}

	return eventHandler.Handle()
}

func resolve(e *EventPayload) EventHandler {
	switch e.Type {
	case EventInstallation:
		return &InstallationEventHandler{event: e}
	case EventPullRequest:
		return &PullRequestEventHandler{event: e}
	default:
		return nil
	}
}
