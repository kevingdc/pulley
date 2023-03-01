package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kevingdc/pulley/pkg/app"
)

type UserService struct {
	DB *sql.DB
}

func (s *UserService) Create(user *app.User) error {
	var lastInsertId int64

	err := s.DB.QueryRow(`
		INSERT INTO users (repository_id, repository_type, chat_id, chat_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`, user.RepositoryID, user.RepositoryType, user.ChatID, user.ChatType).Scan(&lastInsertId)

	if err != nil {
		return fmt.Errorf("User.Create: %v", err)
	}

	createdUser, _ := s.FindByID(lastInsertId)
	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt

	return nil
}

func (s *UserService) FindByID(id int64) (*app.User, error) {
	row := s.DB.QueryRow(`
		SELECT *
		FROM users
		WHERE id = $1;
	`, id)

	return s.scanRow(row)
}

func (s *UserService) FindOneByRepositoryIDAndType(repositoryID string, repositoryType app.Repo) (*app.User, error) {
	row := s.DB.QueryRow(`
	SELECT *
	FROM users
	WHERE repository_id = $1 AND repository_type = $2;
`, repositoryID, repositoryType)
	return s.scanRow(row)
}

func (s *UserService) Exists(user *app.User) bool {
	if user.ID != 0 {
		_, err := s.FindByID(user.ID)
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

	row := s.DB.QueryRow(query, values...)
	_, err := s.scanRow(row)
	return err == nil

}

func (s *UserService) scanRow(row *sql.Row) (*app.User, error) {
	var user app.User
	if err := row.Scan(&user.ID, &user.RepositoryID, &user.RepositoryType, &user.ChatID, &user.ChatType, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("getting user: %v", err)
	}

	return &user, nil
}
