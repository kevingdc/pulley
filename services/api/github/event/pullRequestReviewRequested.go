package event

import (
	"fmt"
	"strconv"

	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
)

func (h *PullRequestEventHandler) handleReviewRequested() (EventHandlerResponse, error) {
	id := strconv.FormatInt(h.prEvent.GetRequestedReviewer().GetID(), 10)
	user, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	content := fmt.Sprintf("**Review Requested**\n>>> %s", h.formattedPRText())

	err = messenger.Send(messenger.Message{
		User:    user,
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
