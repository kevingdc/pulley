package event

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
	"golang.org/x/sync/errgroup"
)

type PullRequestActionHandler interface {
	Handle() (EventHandlerResponse, error)
}

const (
	ActionPRReviewRequested Action = "review_requested"
	ActionPRClosed          Action = "closed"
)

type PullRequestEventHandler struct {
	event   *EventPayload
	action  Action
	prEvent *github.PullRequestEvent
	pr      *github.PullRequest
	repo    *github.Repository
}

func NewPREventHandler(event *EventPayload) *PullRequestEventHandler {
	prEvent := event.Payload.(*github.PullRequestEvent)
	action := Action(prEvent.GetAction())
	pr := prEvent.GetPullRequest()
	repo := prEvent.GetRepo()

	return &PullRequestEventHandler{
		event:   event,
		prEvent: prEvent,
		action:  action,
		pr:      pr,
		repo:    repo,
	}
}

func (h *PullRequestEventHandler) Handle() (EventHandlerResponse, error) {
	switch h.action {
	case ActionPRReviewRequested:
		return h.handleReviewRequested()

	case ActionPRClosed:
		return h.handleClosed()

	default:
		event := h.event.Payload.(*github.PullRequestEvent)

		action := event.GetAction()
		fmt.Printf("Action: %s\n", action)

		fmt.Println()

		repo := event.GetRepo()
		fmt.Printf("Repo name: %s\n", repo.GetName())
		fmt.Printf("Repo full name: %s\n", repo.GetFullName())

		fmt.Println()

		pr := event.GetPullRequest()
		fmt.Printf("PR state: %s\n", pr.GetState())
		fmt.Printf("PR title: %s\n", pr.GetTitle())
		fmt.Printf("PR body: %s\n", pr.GetBody())
		fmt.Printf("PR number: %d\n", pr.GetNumber())
		fmt.Printf("PR URL: %s\n", pr.GetURL())
		fmt.Printf("PR HTML URL: %s\n", pr.GetHTMLURL())
		fmt.Printf("PR Assignees: %+v\n", pr.Assignees)
		fmt.Printf("PR Reviewers: %+v\n", pr.RequestedReviewers)

		return nil, nil
	}
}

func (h *PullRequestEventHandler) formattedPRText() string {
	pr := h.pr
	title := fmt.Sprintf("__**#%d %s**__", pr.GetNumber(), pr.GetTitle())
	url := h.pr.GetHTMLURL()

	commitLabel := "commit"
	if pr.GetCommits() > 1 {
		commitLabel = "commits"
	}

	fileLabel := "file"
	if pr.GetChangedFiles() > 1 {
		fileLabel = "files"
	}
	changeDetails := fmt.Sprintf("*%d %s, %d %s changed*", pr.GetCommits(), commitLabel, pr.GetChangedFiles(), fileLabel)

	body := h.pr.GetBody()

	return fmt.Sprintf("%s\n%s\n%s\n\n%s", title, url, changeDetails, body)
}

func (h *PullRequestEventHandler) getRequestedReviewerUsers() []*user.User {
	var reviewers []*user.User

	for _, reviewer := range h.pr.RequestedReviewers {
		id := user.ToRepoID(reviewer.GetID())
		u, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
		if err != nil {
			continue
		}

		reviewers = append(reviewers, u)
	}

	return reviewers
}

func (h *PullRequestEventHandler) getAssigneeUsers() []*user.User {
	var assignees []*user.User

	for _, assignee := range h.pr.Assignees {
		id := user.ToRepoID(assignee.GetID())
		u, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
		if err != nil {
			continue
		}

		assignees = append(assignees, u)
	}

	return assignees
}

func (h *PullRequestEventHandler) messageUsers(u []*user.User, content string) error {
	g := new(errgroup.Group)

	for _, user := range u {
		user := user
		g.Go(func() error {
			return h.messageUser(user, content)
		})
	}

	return g.Wait()
}

func (h *PullRequestEventHandler) messageUser(u *user.User, content string) error {
	err := messenger.Send(messenger.Message{
		User:    u,
		Content: content,
	})
	if err != nil {
		return err
	}

	return nil
}
