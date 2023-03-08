package github

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
)

type PullRequest interface {
	GetSender() *github.User
	GetOwner() *github.User
	GetAssignees() []*github.User
}

type PullRequestAssignee interface {
	GetAssignee() *github.User
}

type PullRequestRequestedReviewer interface {
	GetRequestedReviewer() *github.User
}

type PRUserService struct {
	userService app.UserService
}

func NewPRUserService(userService app.UserService) *PRUserService {
	return &PRUserService{
		userService: userService,
	}
}

func (s *PRUserService) GetAffectedUsers(pr PullRequest) []*app.User {
	users := []*app.User{}

	if !s.isPROwnerSameAsEventSender(pr) {
		user := s.GetOwnerUser(pr)
		users = append(users, user)
	}

	senderID := pr.GetSender().GetID()
	for _, assignee := range pr.GetAssignees() {
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

func (s *PRUserService) GetOwnerUser(pr PullRequest) *app.User {
	ownerID := pr.GetOwner().GetID()
	return s.FindUserByRepoID(ownerID)
}

func (s *PRUserService) GetAssigneeUser(pr PullRequestAssignee) *app.User {
	assigneeID := pr.GetAssignee().GetID()
	return s.FindUserByRepoID(assigneeID)
}

func (s *PRUserService) GetRequestedReviewerUser(pr PullRequestRequestedReviewer) *app.User {
	reviewerID := pr.GetRequestedReviewer().GetID()
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

func (s *PRUserService) IsUserSameAsSender(u *github.User, pr PullRequest) bool {
	return u.GetID() == pr.GetSender().GetID()
}

func (s *PRUserService) isPROwnerSameAsEventSender(pr PullRequest) bool {
	return s.IsUserSameAsSender(pr.GetOwner(), pr)
}
