package pr

import (
	"github.com/kevingdc/pulley/pkg/app"
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

	err := h.handler.messageUsers(usersToMessage, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ClosedReopenedActionHandler) generateMessageContent() *app.MessageContent {
	var (
		actionLabel string
		color       app.Color
	)

	switch {
	case h.handler.action == event.ActionPRReopened:
		actionLabel = "Reopened"
		color = app.ColorBlue

	case h.handler.pr.GetMerged():
		actionLabel = "Merged"
		color = app.ColorPurple

	default:
		actionLabel = "Closed"
		color = app.ColorRed
	}

	return h.handler.generateMessageContent(actionLabel, color)
}
