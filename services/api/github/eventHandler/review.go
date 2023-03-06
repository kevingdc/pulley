package eventhandler

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type PullRequestReviewEventHandler struct {
	event       *event.Payload
	app         *app.App
	userService app.UserService
	reviewEvent *github.PullRequestReviewEvent
	pr          *github.PullRequest
	prReview    *github.PullRequestReview
	repo        *github.Repository
}

func NewPullRequestReviewEventHandler(e *event.Payload) *PullRequestReviewEventHandler {
	reviewEvent := e.Payload.(*github.PullRequestReviewEvent)
	pr := reviewEvent.GetPullRequest()
	repo := reviewEvent.GetRepo()
	prReview := reviewEvent.GetReview()

	return &PullRequestReviewEventHandler{
		event:       e,
		app:         e.App,
		userService: e.App.UserService,
		reviewEvent: reviewEvent,
		pr:          pr,
		prReview:    prReview,
		repo:        repo,
	}
}

func (h *PullRequestReviewEventHandler) Handle() (event.HandlerResponse, error) {
	// usersToMessage :=

	return nil, nil
}
