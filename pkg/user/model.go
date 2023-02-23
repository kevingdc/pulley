package user

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
