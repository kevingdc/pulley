package pr

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type ClosedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *ClosedActionHandler) Handle() (event.HandlerResponse, error) {
	usersToMessage := h.usersToMessage()
	if len(usersToMessage) == 0 {
		return nil, nil
	}

	action := "Closed"
	if h.handler.pr.GetMerged() {
		action = "Merged"
	}
	closerUser := h.handler.prEvent.GetSender().GetLogin()

	content := fmt.Sprintf("**Pull Request %s** *by %s*\n>>> %s", action, closerUser, h.handler.formattedPRText())

	err := h.handler.messageUsers(usersToMessage, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ClosedActionHandler) usersToMessage() []*app.User {
	userService := h.handler.app.UserService
	eventSender := h.handler.prEvent.GetSender()
	prOwner := h.handler.pr.GetUser()

	usersToMessage := []*app.User{}

	if didOwnerClosePR := eventSender.GetID() == prOwner.GetID(); !didOwnerClosePR {
		id := idconv.ToRepoID(prOwner.GetID())
		user, err := userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
		if err == nil {
			usersToMessage = append(usersToMessage, user)
		}
	}

	usersToMessage = append(usersToMessage, h.handler.getAssigneeUsers()...)

	return usersToMessage
}
