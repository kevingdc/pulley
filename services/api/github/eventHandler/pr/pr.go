package pr

import (
	"fmt"

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

func (h *PullRequestEventHandler) generateMessageContent(actionLabel string, color app.Color) *app.MessageContent {
	actingUser := h.eventSender()
	pr := h.pr

	commitLabel := "commit"
	if pr.GetCommits() > 1 {
		commitLabel = "commits"
	}

	fileLabel := "file"
	if pr.GetChangedFiles() > 1 {
		fileLabel = "files"
	}

	return &app.MessageContent{
		URL:       pr.GetHTMLURL(),
		Title:     fmt.Sprintf("#%d %s", pr.GetNumber(), pr.GetTitle()),
		Subtitle:  fmt.Sprintf("*%d %s, %d %s changed*", pr.GetCommits(), commitLabel, pr.GetChangedFiles(), fileLabel),
		Body:      pr.GetBody(),
		Color:     color,
		Thumbnail: pr.GetUser().GetAvatarURL(),
		Author: &app.MessageAuthor{
			Name:      actingUser.GetLogin(),
			URL:       actingUser.GetHTMLURL(),
			AvatarURL: actingUser.GetAvatarURL(),
		},
		Header: fmt.Sprintf("PR %s", actionLabel),
		Footer: h.repo.GetFullName(),
	}
}

func (h *PullRequestEventHandler) eventSender() *github.User {
	return h.prEvent.GetSender()
}
