package app

import "time"

type Repo string

const (
	RepoGitHub Repo = "github"
)

type Chat string

const (
	ChatDiscord Chat = "discord"
)

type User struct {
	ID             int64
	RepositoryID   string
	RepositoryType Repo
	ChatID         string
	ChatType       Chat
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ChatConfig struct {
	ChatID   string `json:"chatId"`
	ChatType string `json:"chatType"`
}

type UserService interface {
	Create(user *User) error
	FindByID(id int64) (*User, error)
	FindByRepositoryIDAndType(repositoryID string, repositoryType Repo) ([]User, error)
	FindOneByRepositoryIDAndType(repositoryID string, repositoryType Repo) (*User, error)
	FindByChatIDAndType(chatID string, chatType Chat) ([]User, error)
	FindOneByChatIDAndType(chatID string, chatType Chat) (*User, error)
	Exists(user *User) bool
}
