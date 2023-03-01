package pr

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type ReviewRequestedActionHandler struct {
	handler *PullRequestEventHandler
}

func (h *ReviewRequestedActionHandler) Handle() (event.HandlerResponse, error) {
	user, err := h.userToMessage()
	if err != nil {
		return nil, nil
	}

	content := fmt.Sprintf("**Review Requested**\n>>> %s", h.handler.formattedPRText())

	err = h.handler.messageUser(user, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ReviewRequestedActionHandler) userToMessage() (*app.User, error) {
	userService := h.handler.userService

	id := idconv.ToRepoID(h.handler.prEvent.GetRequestedReviewer().GetID())
	user, err := userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
	if err != nil {
		return nil, err
	}

	return user, err
}
