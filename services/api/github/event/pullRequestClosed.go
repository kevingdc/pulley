package event

import (
	"fmt"
	"strconv"

	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
)

func (h *PullRequestEventHandler) handleClosed() (EventHandlerResponse, error) {
	eventSender := h.prEvent.GetSender()
	prOwner := h.pr.GetUser()

	if didOwnerClosePR := eventSender.GetID() == prOwner.GetID(); didOwnerClosePR {
		return nil, nil
	}

	id := strconv.FormatInt(prOwner.GetID(), 10)
	user, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	closerUser := eventSender.GetLogin()

	action := "Closed"
	if h.pr.GetMerged() {
		action = "Merged"
	}

	content := fmt.Sprintf("**Pull Request %s** *by %s*\n>>> %s", action, closerUser, h.formattedPRText())

	err = messenger.Send(messenger.Message{
		User:    user,
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
