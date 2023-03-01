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
	user := h.userToMessage()
	if user == nil {
		return nil, nil
	}

	content := fmt.Sprintf("**Review Requested**\n>>> %s", h.handler.formattedPRText())

	err := h.handler.messageUser(user, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ReviewRequestedActionHandler) userToMessage() *app.User {
	userService := h.handler.userService

	requestedReviewerID := h.handler.prEvent.GetRequestedReviewer().GetID()

	if requestedReviewerID == h.handler.eventSender().GetID() {
		return nil
	}

	repoID := idconv.ToRepoID(requestedReviewerID)
	user, _ := userService.FindOneByRepositoryIDAndType(repoID, app.RepoGitHub)

	return user
}
