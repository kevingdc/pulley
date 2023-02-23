-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  repository_id VARCHAR (20) NOT NULL,
  repository_type VARCHAR (20) NOT NULL,
  chat_id VARCHAR (20) NOT NULL,
  chat_type VARCHAR (20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
