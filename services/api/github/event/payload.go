package event

import (
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/config"
)

type EventType string

const (
	EventPullRequest              EventType = "pull_request"
	EventPullRequestReview        EventType = "pull_request_review"
	EventPullRequestReviewComment EventType = "pull_request_review_comment"
	EventInstallation             EventType = "installation"
)

type EventPayload struct {
	Config  *config.Config
	Payload interface{}
	Github  *github.Client
	Type    EventType
}
