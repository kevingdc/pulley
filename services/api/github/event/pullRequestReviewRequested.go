package event

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/user"
)

func (h *PullRequestEventHandler) handleReviewRequested() (EventHandlerResponse, error) {
	id := user.ToRepoID(h.prEvent.GetRequestedReviewer().GetID())
	user, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	content := fmt.Sprintf("**Review Requested**\n>>> %s", h.formattedPRText())

	err = h.messageUser(user, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
