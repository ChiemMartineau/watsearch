-- +goose Up
-- +goose StatementBegin
CREATE TABLE courses (
  id   TEXT PRIMARY KEY,
  title TEXT,
  description TEXT,
  credits TEXT,
  prerequisites TEXT,
  antirequisites TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE courses;
-- +goose StatementEnd
