package event

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
)

type Payload struct {
	App     *app.App
	Payload interface{}
	Github  *github.Client
	Type    WebhookType
}

type HandlerResponse interface{}

type Handler interface {
	Handle() (HandlerResponse, error)
}

type WebhookType string

const (
	TypeInstallation             WebhookType = "installation"
	TypePullRequest              WebhookType = "pull_request"
	TypePullRequestReview        WebhookType = "pull_request_review"
	TypePullRequestReviewComment WebhookType = "pull_request_review_comment"
)

type Action string

const (
	ActionInstalled   Action = "created"
	ActionUninstalled Action = "deleted"

	ActionPRReviewRequested Action = "review_requested"
	ActionPRClosed          Action = "closed"
	ActionPRReopened        Action = "reopened"
	ActionPRAssigned        Action = "assigned"
	ActionPRUnassigned      Action = "unassigned"
)

type ReviewState string

const (
	ReviewDismissed        ReviewState = "dismissed"
	ReviewApproved         ReviewState = "approved"
	ReviewCommented        ReviewState = "commented"
	ReviewChangesRequested ReviewState = "changes_requested"
)
