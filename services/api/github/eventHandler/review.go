package eventhandler

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	internalgithub "github.com/kevingdc/pulley/pkg/github"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type PullRequestReviewEventHandler struct {
	payload       *event.Payload
	app           *app.App
	prUserService *internalgithub.PRUserService
	reviewEvent   *github.PullRequestReviewEvent
	pr            *github.PullRequest
	prReview      *github.PullRequestReview
	repo          *github.Repository
}

func NewPullRequestReviewEventHandler(e *event.Payload) *PullRequestReviewEventHandler {
	reviewEvent := e.Payload.(*github.PullRequestReviewEvent)
	pr := reviewEvent.GetPullRequest()
	repo := reviewEvent.GetRepo()
	prReview := reviewEvent.GetReview()

	return &PullRequestReviewEventHandler{
		payload:       e,
		app:           e.App,
		prUserService: internalgithub.NewPRUserService(e.App.UserService),
		reviewEvent:   reviewEvent,
		pr:            pr,
		prReview:      prReview,
		repo:          repo,
	}
}

func (h *PullRequestReviewEventHandler) Handle() (event.HandlerResponse, error) {
	affectedUsers := h.prUserService.GetAffectedUsers(h.reviewEvent)
	if len(affectedUsers) == 0 {
		return nil, nil
	}

	err := messenger.SendToUsers(affectedUsers, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *PullRequestReviewEventHandler) generateMessageContent() *app.MessageContent {
	builder := internalgithub.NewPRMessageBuilder()
	builder.SetPR(h.pr)
	builder.SetPRReview(h.prReview)
	builder.SetSender(h.reviewEvent.GetSender())
	builder.SetEvent(h.event())
	builder.SetMessageType(internalgithub.TypePRReview)

	return builder.Build()
}

func (h *PullRequestReviewEventHandler) event() event.Event {
	switch event.ReviewState(h.prReview.GetState()) {
	case event.ReviewApproved:
		return &event.PRApproved{}

	case event.ReviewCommented:
		return &event.PRCommented{}

	case event.ReviewChangesRequested:
		return &event.PRChangesRequested{}

	default:
		return nil
	}
}
