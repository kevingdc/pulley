package event

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/user"
)

func (h *PullRequestEventHandler) handleClosed() (EventHandlerResponse, error) {
	eventSender := h.prEvent.GetSender()
	prOwner := h.pr.GetUser()

	usersToMessage := []*user.User{}

	if didOwnerClosePR := eventSender.GetID() == prOwner.GetID(); !didOwnerClosePR {
		id := user.ToRepoID(prOwner.GetID())
		user, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
		if err != nil {
			return nil, err
		}

		usersToMessage = append(usersToMessage, user)
	}

	usersToMessage = append(usersToMessage, h.getAssigneeUsers()...)

	action := "Closed"
	if h.pr.GetMerged() {
		action = "Merged"
	}
	closerUser := eventSender.GetLogin()

	content := fmt.Sprintf("**Pull Request %s** *by %s*\n>>> %s", action, closerUser, h.formattedPRText())

	h.messageUsers(usersToMessage, content)

	return nil, nil
}
