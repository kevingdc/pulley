package event

import "fmt"

type PullRequestEventHandler struct {
	event *EventPayload
}

func (h *PullRequestEventHandler) Handle() (EventHandlerResponse, error) {
	fmt.Println("PullRequestEventHandler.Handle()")
	fmt.Printf("%+v\n", h.event)
	return nil, nil
}
