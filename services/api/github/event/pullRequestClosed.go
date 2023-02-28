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
	user, err := user.GetByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	closerUser := eventSender.GetLogin()

	content := fmt.Sprintf("**Pull Request Closed** *by %s*\n>>> %s", closerUser, h.formattedPRText())

	messenger.Send(messenger.Message{
		User:    user,
		Content: content,
	})

	return nil, nil
}
