package pr

import "github.com/google/go-github/v50/github"

type prAdapter struct {
	prEvent *github.PullRequestEvent
}

func (a *prAdapter) GetSender() *github.User {
	return a.prEvent.GetSender()
}

func (a *prAdapter) GetOwner() *github.User {
	return a.prEvent.GetPullRequest().GetUser()
}

func (a *prAdapter) GetAssignees() []*github.User {
	return a.prEvent.GetPullRequest().Assignees
}

func (a *prAdapter) GetAssignee() *github.User {
	return a.prEvent.GetAssignee()
}

func (a *prAdapter) GetRequestedReviewer() *github.User {
	return a.prEvent.GetRequestedReviewer()
}
