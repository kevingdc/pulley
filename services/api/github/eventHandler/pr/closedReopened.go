package pr

import (
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type ClosedReopenedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *ClosedReopenedActionHandler) Handle() (event.HandlerResponse, error) {
	affectedUsers := h.handler.prUserService.GetAffectedUsers(h.handler.adapter)
	if len(affectedUsers) == 0 {
		return nil, nil
	}

	err := messenger.SendToUsers(affectedUsers, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ClosedReopenedActionHandler) generateMessageContent() *app.MessageContent {
	var e event.Event

	switch {
	case h.handler.action == event.ActionPRReopened:
		e = &event.PRReopened{}

	case h.handler.pr.GetMerged():
		e = &event.PRMerged{}

	default:
		e = &event.PRClosed{}
	}

	return h.handler.generateMessageContent(e)
}
