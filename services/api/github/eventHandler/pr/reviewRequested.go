package pr

import (
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type ReviewRequestedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *ReviewRequestedActionHandler) Handle() (event.HandlerResponse, error) {
	user := h.userToMessage()
	if user == nil {
		return nil, nil
	}

	err := messenger.SendToUser(user, h.handler.generateMessageContent("Review Requested", app.ColorYellow))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ReviewRequestedActionHandler) userToMessage() *app.User {
	prUserService := h.handler.prUserService

	requestedReviewer := h.handler.prEvent.GetRequestedReviewer()
	if prUserService.IsUserSameAsSender(requestedReviewer, h.handler.prEvent) {
		return nil
	}

	return prUserService.GetRequestedReviewerUser(h.handler.prEvent)
}
