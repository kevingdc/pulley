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

	createdUser, _ := FindByID(lastInsertId)
	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt

	return lastInsertId, nil
}

func (user *User) Exists() bool {
	if user.ID != 0 {
		_, err := FindByID(user.ID)
		return err == nil
	}

	var (
		columns []string
		values  []interface{}
	)

	if user.RepositoryID != "" {
		columns = append(columns, "repository_id")
		values = append(values, user.RepositoryID)
	}

	if user.RepositoryType != "" {
		columns = append(columns, "repository_type")
		values = append(values, user.RepositoryType)
	}

	if user.ChatID != "" {
		columns = append(columns, "chat_id")
		values = append(values, user.ChatID)
	}

	if user.ChatType != "" {
		columns = append(columns, "chat_type")
		values = append(values, user.ChatType)
	}

	if len(columns) == 0 {
		return false
	}

	var whereConditions string
	for i, column := range columns {
		if i == 0 {
			whereConditions = fmt.Sprintf("%s = $%d", column, i+1)
			continue
		}

		whereConditions = fmt.Sprintf("%s AND %s = $%d", whereConditions, column, i+1)
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM users
		WHERE %s;
	`, whereConditions)

	row := db.DB.QueryRow(query, values...)
	_, err := scanRow(row)
	return err == nil
}

func FindByID(id int64) (*User, error) {
	row := db.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE id = $1;
	`, id)

	return scanRow(row)
}

func FindOneByRepositoryIDAndType(repositoryID string, repositoryType Repo) (*User, error) {
	row := db.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE repository_id = $1 AND repository_type = $2;
	`, repositoryID, repositoryType)
	return scanRow(row)
}

func scanRow(row *sql.Row) (*User, error) {
	var user User
	if err := row.Scan(&user.ID, &user.RepositoryID, &user.RepositoryType, &user.ChatID, &user.ChatType, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("getting user: %v", err)
	}

	return &user, nil
}
