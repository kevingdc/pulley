package user

// import (
// 	"gorm.io/gorm"
// )

type Repo string

const (
	RepoGitHub Repo = "github"
)

type Chat string

const (
	ChatDiscord Chat = "discord"
)

type User struct {
	// gorm.Model
	RepositoryID   string
	RepositoryType Repo
	ChatID         string
	ChatType       Chat
}
