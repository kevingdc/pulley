package pr

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	internalgithub "github.com/kevingdc/pulley/pkg/github"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type PullRequestEventHandler struct {
	event         *event.Payload
	app           *app.App
	prUserService *internalgithub.PRUserService
	action        event.Action
	prEvent       *github.PullRequestEvent
	pr            *github.PullRequest
	repo          *github.Repository
	adapter       *prAdapter
}

func New(e *event.Payload) *PullRequestEventHandler {
	prEvent := e.Payload.(*github.PullRequestEvent)
	action := event.Action(prEvent.GetAction())
	pr := prEvent.GetPullRequest()
	repo := prEvent.GetRepo()
	prUserService := internalgithub.NewPRUserService(e.App.UserService)

	return &PullRequestEventHandler{
		event:         e,
		prEvent:       prEvent,
		action:        action,
		pr:            pr,
		repo:          repo,
		app:           e.App,
		prUserService: prUserService,
		adapter:       &prAdapter{prEvent: prEvent},
	}
}

func (h *PullRequestEventHandler) Handle() (event.HandlerResponse, error) {
	handler := h.resolve()

	if handler == nil {
		return nil, nil
	}

	return handler.Handle()
}

func (h *PullRequestEventHandler) resolve() event.Handler {
	switch h.action {
	case event.ActionPRReviewRequested:
		return &ReviewRequestedActionHandler{handler: h}

	case event.ActionPRClosed:
		return &ClosedReopenedActionHandler{handler: h}

	case event.ActionPRReopened:
		return &ClosedReopenedActionHandler{handler: h}

	case event.ActionPRAssigned:
		return &AssignedUnassignedActionHandler{handler: h}

	case event.ActionPRUnassigned:
		return &AssignedUnassignedActionHandler{handler: h}

	default:
		return nil
	}
}

func (h *PullRequestEventHandler) generateMessageContent(event event.Event) *app.MessageContent {
	builder := internalgithub.NewPRMessageBuilder()
	builder.SetPR(h.pr)
	builder.SetSender(h.prEvent.GetSender())
	builder.SetEvent(event)
	builder.SetMessageType(internalgithub.TypePRMessage)

	return builder.Build()
}
