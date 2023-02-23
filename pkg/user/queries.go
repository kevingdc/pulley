package user

import (
	"database/sql"
	"fmt"

	"github.com/kevingdc/pulley/pkg/db"
)

func (user *User) Create() (int64, error) {
	var lastInsertId int64

	err := db.DB.QueryRow(`
		INSERT INTO users (repository_id, repository_type, chat_id, chat_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, user.RepositoryID, user.RepositoryType, user.ChatID, user.ChatType).Scan(&lastInsertId)

	if err != nil {
		return 0, fmt.Errorf("User.Create: %v", err)
	}

	createdUser, _ := GetByID(lastInsertId)
	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt

	return lastInsertId, nil
}

func GetByID(id int64) (User, error) {
	var user User

	row := db.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE id = $1;
	`, id)
	if err := row.Scan(&user.ID, &user.RepositoryID, &user.RepositoryType, &user.ChatID, &user.ChatType, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("getByID(%d): user not found", id)
		}
		return user, fmt.Errorf("getByID(%d): %v", id, err)
	}

	return user, nil
}

func GetByRepositoryIDAndType(repositoryID string, repositoryType Repo) (User, error) {
	var user User

	row := db.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE repository_id = $1 AND repository_type = $2;
	`, repositoryID, repositoryType)
	if err := row.Scan(&user.ID, &user.RepositoryID, &user.RepositoryType, &user.ChatID, &user.ChatType, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("getByRepositoryIDAndType(%q, %q): user not found", repositoryID, repositoryType)
		}
		return user, fmt.Errorf("getByRepositoryIDAndType(%q, %q): %v", repositoryID, repositoryType, err)
	}
	return user, nil
}
