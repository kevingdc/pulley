package pr

import (
	"github.com/kevingdc/pulley/services/api/github/event"
)

type ClosedReopenedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *ClosedReopenedActionHandler) Handle() (event.HandlerResponse, error) {
	usersToMessage := h.handler.affectedUsers()
	if len(usersToMessage) == 0 {
		return nil, nil
	}

	err := h.handler.messageUsers(usersToMessage, h.handler.generateMessageContent(h.actionLabel()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ClosedReopenedActionHandler) actionLabel() string {
	switch {
	case h.handler.action == event.ActionPRReopened:
		return "Reopened"
	case h.handler.pr.GetMerged():
		return "Merged"
	default:
		return "Closed"
	}
}
