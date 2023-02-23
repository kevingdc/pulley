package event

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
)

type PullRequestEventHandler struct {
	event *EventPayload
}

func (h *PullRequestEventHandler) Handle() (EventHandlerResponse, error) {
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

	id := strconv.FormatInt(event.GetRequestedReviewer().GetID(), 10)
	user, err := user.GetByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	messenger.Send(messenger.Message{
		User:    user,
		Content: fmt.Sprintf("Hey, you have a new pull request to review: %s", pr.GetHTMLURL()),
	})

	return nil, nil
}
