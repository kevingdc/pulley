package pr

import (
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
	"golang.org/x/sync/errgroup"
)

type PullRequestEventHandler struct {
	event       *event.Payload
	action      event.Action
	prEvent     *github.PullRequestEvent
	pr          *github.PullRequest
	repo        *github.Repository
	app         *app.App
	userService app.UserService
}

func New(e *event.Payload) *PullRequestEventHandler {
	prEvent := e.Payload.(*github.PullRequestEvent)
	action := event.Action(prEvent.GetAction())
	pr := prEvent.GetPullRequest()
	repo := prEvent.GetRepo()

	return &PullRequestEventHandler{
		event:       e,
		prEvent:     prEvent,
		action:      action,
		pr:          pr,
		repo:        repo,
		app:         e.App,
		userService: e.App.UserService,
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

func (h *PullRequestEventHandler) eventSender() *github.User {
	return h.prEvent.GetSender()
}

func (h *PullRequestEventHandler) prOwner() *github.User {
	return h.pr.GetUser()
}

func (h *PullRequestEventHandler) prOwnerUser() *app.User {
	id := idconv.ToRepoID(h.prOwner().GetID())
	u, err := h.userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
	if err != nil {
		return nil
	}

	return u
}

func (h *PullRequestEventHandler) isPROwnerSameAsEventSender() bool {
	return h.prOwner().GetID() == h.eventSender().GetID()
}

func (h *PullRequestEventHandler) generateMessageContent(action string) string {
	actingUser := h.eventSender().GetLogin()
	prText := h.formattedPRText()
	return fmt.Sprintf("**Pull Request %s** *by %s*\n>>> %s", action, actingUser, prText)
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

func (h *PullRequestEventHandler) requestedReviewerUsers() []*app.User {
	var reviewers []*app.User

	for _, reviewer := range h.pr.RequestedReviewers {
		id := idconv.ToRepoID(reviewer.GetID())
		u, err := h.userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
		if err != nil {
			continue
		}

		reviewers = append(reviewers, u)
	}

	return reviewers
}

func (h *PullRequestEventHandler) assigneeUsers() []*app.User {
	var assignees []*app.User

	for _, assignee := range h.pr.Assignees {
		id := idconv.ToRepoID(assignee.GetID())
		u, err := h.userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
		if err != nil {
			continue
		}

		assignees = append(assignees, u)
	}

	return assignees
}

func (h *PullRequestEventHandler) affectedUsers() []*app.User {
	users := []*app.User{}

	if !h.isPROwnerSameAsEventSender() {
		user := h.prOwnerUser()
		users = append(users, user)
	}

	senderID := h.eventSender().GetID()
	for _, assignee := range h.pr.Assignees {
		assigneeID := assignee.GetID()
		if assigneeID == senderID {
			continue
		}

		repoID := idconv.ToRepoID(assigneeID)
		u, err := h.userService.FindOneByRepositoryIDAndType(repoID, app.RepoGitHub)
		if err != nil {
			continue
		}

		users = append(users, u)
	}

	return users
}

func (h *PullRequestEventHandler) messageUsers(u []*app.User, content string) error {
	g := new(errgroup.Group)

	for _, user := range u {
		if user == nil {
			continue
		}

		user := user
		g.Go(func() error {
			return h.messageUser(user, content)
		})
	}

	return g.Wait()
}

func (h *PullRequestEventHandler) messageUser(u *app.User, content string) error {
	err := messenger.Send(messenger.Message{
		User:    u,
		Content: content,
	})
	if err != nil {
		return err
	}

	return nil
}
