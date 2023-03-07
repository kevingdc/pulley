package github

import (
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type PRMessageType string

const (
	TypePRMessage PRMessageType = "PRMessage"
	TypePRReview  PRMessageType = "PRReview"
)

type PRMessageBuilder struct {
	pr          *github.PullRequest
	prOwner     *github.User
	prReview    *github.PullRequestReview
	repo        *github.Repository
	sender      *github.User
	event       event.Event
	messageType PRMessageType
}

func NewPRMessageBuilder() *PRMessageBuilder {
	return &PRMessageBuilder{}
}

func (b *PRMessageBuilder) SetPR(pr *github.PullRequest) *PRMessageBuilder {
	b.pr = pr
	b.prOwner = pr.GetUser()
	b.repo = pr.GetBase().GetRepo()
	return b
}

func (b *PRMessageBuilder) SetPRReview(prReview *github.PullRequestReview) *PRMessageBuilder {
	b.prReview = prReview
	return b
}

func (b *PRMessageBuilder) SetSender(sender *github.User) *PRMessageBuilder {
	b.sender = sender
	return b
}

func (b *PRMessageBuilder) SetEvent(e event.Event) *PRMessageBuilder {
	b.event = e
	return b
}

func (b *PRMessageBuilder) SetMessageType(messageType PRMessageType) *PRMessageBuilder {
	b.messageType = messageType
	return b
}

func (b *PRMessageBuilder) Build() *app.MessageContent {
	switch b.messageType {
	case TypePRMessage:
		return b.prMessageContent()
	case TypePRReview:
		return b.prReviewMessageContent()
	default:
		return nil
	}
}

func (b *PRMessageBuilder) prMessageContent() *app.MessageContent {
	if b.pr == nil {
		return nil
	}

	commitLabel := "commit"
	if b.pr.GetCommits() > 1 {
		commitLabel = "commits"
	}

	fileLabel := "file"
	if b.pr.GetChangedFiles() > 1 {
		fileLabel = "files"
	}

	baseMessage := b.baseMessage()
	baseMessage.Body = b.pr.GetBody()
	baseMessage.Subtitle = fmt.Sprintf("*%d %s, %d %s changed*", b.pr.GetCommits(), commitLabel, b.pr.GetChangedFiles(), fileLabel)
	baseMessage.Thumbnail = b.pr.GetUser().GetAvatarURL()

	return baseMessage
}

func (b *PRMessageBuilder) prReviewMessageContent() *app.MessageContent {
	if b.pr == nil || b.prReview == nil {
		return nil
	}

	baseMessage := b.baseMessage()
	baseMessage.Body = b.prReview.GetBody()
	baseMessage.Thumbnail = b.prReview.GetUser().GetAvatarURL()

	return baseMessage
}

func (b *PRMessageBuilder) baseMessage() *app.MessageContent {
	return &app.MessageContent{
		URL:    b.pr.GetHTMLURL(),
		Title:  b.messageTitle(),
		Color:  b.event.Color(),
		Author: b.messageAuthor(),
		Header: fmt.Sprintf("PR %s", b.event.Name()),
		Footer: b.repo.GetFullName(),
	}
}

func (b *PRMessageBuilder) messageTitle() string {
	return fmt.Sprintf("#%d %s", b.pr.GetNumber(), b.pr.GetTitle())
}

func (b *PRMessageBuilder) messageAuthor() *app.MessageAuthor {
	if b.sender == nil {
		return nil
	}

	return &app.MessageAuthor{
		Name:      b.sender.GetLogin(),
		URL:       b.sender.GetHTMLURL(),
		AvatarURL: b.sender.GetAvatarURL(),
	}
}
