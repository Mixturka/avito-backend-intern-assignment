-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (team_name TEXT PRIMARY KEY);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd