package github

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
)

type PullRequestEvent interface {
	GetPullRequest() *github.PullRequest
	GetSender() *github.User
}

type PRUserService struct {
	userService app.UserService
}

func NewPRUserService(userService app.UserService) *PRUserService {
	return &PRUserService{
		userService: userService,
	}
}

func (s *PRUserService) IsUserSameAsSender(u *github.User, e PullRequestEvent) bool {
	return u.GetID() == e.GetSender().GetID()
}

func (s *PRUserService) GetAffectedUsers(e PullRequestEvent) []*app.User {
	pr := e.GetPullRequest()
	users := []*app.User{}

	if !s.isPROwnerSameAsEventSender(e) {
		user := s.GetOwnerUserFromPR(pr)
		users = append(users, user)
	}

	senderID := e.GetSender().GetID()
	for _, assignee := range pr.Assignees {
		assigneeID := assignee.GetID()
		if assigneeID == senderID {
			continue
		}

		u := s.FindUserByRepoID(assigneeID)
		if u == nil {
			continue
		}

		users = append(users, u)
	}

	return users
}

func (s *PRUserService) GetOwnerUserFromEvent(e PullRequestEvent) *app.User {
	return s.GetOwnerUserFromPR(e.GetPullRequest())
}

func (s *PRUserService) GetOwnerUserFromPR(pr *github.PullRequest) *app.User {
	ownerID := pr.GetUser().GetID()
	return s.FindUserByRepoID(ownerID)
}

func (s *PRUserService) GetAssigneeUser(e *github.PullRequestEvent) *app.User {
	assigneeID := e.GetAssignee().GetID()
	return s.FindUserByRepoID(assigneeID)
}

func (s *PRUserService) GetRequestedReviewerUser(e *github.PullRequestEvent) *app.User {
	reviewerID := e.GetRequestedReviewer().GetID()
	return s.FindUserByRepoID(reviewerID)
}

func (s *PRUserService) FindUserByRepoID(repoID int64) *app.User {
	convertedID := idconv.ToRepoID(repoID)
	u, err := s.userService.FindOneByRepositoryIDAndType(convertedID, app.RepoGitHub)
	if err != nil {
		return nil
	}

	return u
}

func (s *PRUserService) isPROwnerSameAsEventSender(e PullRequestEvent) bool {
	pr := e.GetPullRequest()
	return pr.GetUser().GetID() == e.GetSender().GetID()
}
