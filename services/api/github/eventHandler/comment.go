package eventhandler

import (
	"context"

	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	internalgithub "github.com/kevingdc/pulley/pkg/github"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type PullRequestCommentEventHandler struct {
	payload       *event.Payload
	app           *app.App
	prUserService *internalgithub.PRUserService
	commentEvent  *github.IssueCommentEvent
	issue         *github.Issue
	comment       *github.IssueComment
}

func NewPullRequestCommentEventHandler(e *event.Payload) *PullRequestCommentEventHandler {
	commentEvent := e.Payload.(*github.IssueCommentEvent)

	return &PullRequestCommentEventHandler{
		payload:       e,
		app:           e.App,
		prUserService: internalgithub.NewPRUserService(e.App.UserService),
		commentEvent:  commentEvent,
		issue:         commentEvent.GetIssue(),
		comment:       commentEvent.GetComment(),
	}
}

func (h *PullRequestCommentEventHandler) Handle() (event.HandlerResponse, error) {
	if !h.issue.IsPullRequest() {
		return nil, nil
	}

	affectedUsers := h.prUserService.GetAffectedUsers(h)
	if len(affectedUsers) == 0 {
		return nil, nil
	}

	err := messenger.SendToUsers(affectedUsers, h.generateMessageContent())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *PullRequestCommentEventHandler) GetSender() *github.User {
	return h.commentEvent.GetSender()
}

func (h *PullRequestCommentEventHandler) GetOwner() *github.User {
	return h.issue.GetUser()
}

func (h *PullRequestCommentEventHandler) GetAssignees() []*github.User {
	return h.issue.Assignees
}

func (h *PullRequestCommentEventHandler) generateMessageContent() *app.MessageContent {
	client := h.payload.Github
	repo := h.commentEvent.GetRepo()

	pr, _, err := client.PullRequests.Get(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), h.issue.GetNumber())
	if err != nil {
		return nil
	}

	builder := internalgithub.NewPRMessageBuilder()
	builder.SetPR(pr)
	builder.SetPRReview(h.comment)
	builder.SetSender(h.commentEvent.GetSender())
	builder.SetEvent(&event.PRCommented{})
	builder.SetMessageType(internalgithub.TypePRReview)

	return builder.Build()
}
