package pr

import (
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
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

	err := h.handler.messageUser(user, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *AssignedUnassignedActionHandler) userToMessage() *app.User {
	userService := h.handler.userService

	assigneeID := h.handler.prEvent.GetAssignee().GetID()

	if assigneeID == h.handler.eventSender().GetID() {
		return nil
	}

	repoID := idconv.ToRepoID(assigneeID)
	user, _ := userService.FindOneByRepositoryIDAndType(repoID, app.RepoGitHub)

	return user
}

func (h *AssignedUnassignedActionHandler) generateMessageContent() *app.MessageContent {
	var (
		actionLabel string
		color       app.Color
	)

	switch h.handler.action {
	case event.ActionPRAssigned:
		actionLabel = "Assigned"
		color = app.ColorCyan

	case event.ActionPRUnassigned:
		actionLabel = "Unassigned"
		color = app.ColorGrey

	default:
		return nil
	}

	return h.handler.generateMessageContent(actionLabel, color)
}

func (h *AssignedUnassignedActionHandler) actionLabel() string {
	switch h.handler.action {
	case event.ActionPRAssigned:
		return "Assigned"
	case event.ActionPRUnassigned:
		return "Unassigned"
	default:
		return ""
	}
}
