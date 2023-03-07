package pr

import (
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type AssignedUnassignedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *AssignedUnassignedActionHandler) Handle() (event.HandlerResponse, error) {
	user := h.userToMessage()
	if user == nil {
		return nil, nil
	}

	err := messenger.SendToUser(user, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *AssignedUnassignedActionHandler) userToMessage() *app.User {
	prUserService := h.handler.prUserService

	assignee := h.handler.prEvent.GetAssignee()
	if prUserService.IsUserSameAsSender(assignee, h.handler.prEvent) {
		return nil
	}

	return prUserService.GetAssigneeUser(h.handler.prEvent)
}

func (h *AssignedUnassignedActionHandler) generateMessageContent() *app.MessageContent {
	var e event.Event

	switch h.handler.action {
	case event.ActionPRAssigned:
		e = &event.PRAssigned{}

	case event.ActionPRUnassigned:
		e = &event.PRUnassigned{}

	default:
		return nil
	}

	return h.handler.generateMessageContent(e)
}
