-- +goose Up

-- Create enum type if not exists
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
-- ignore error if type already exists
DO $$ BEGIN END; $$;

-- Create pullrequests table
CREATE TABLE IF NOT EXISTS pullrequests (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users (id),
    status pr_status NOT NULL,
    created_at TIMESTAMPTZ,
    merged_at TIMESTAMPTZ
);

-- Create reviewers table
CREATE TABLE IF NOT EXISTS assigned_pr_reviewers (
    pr_id TEXT REFERENCES pullrequests (id) ON DELETE CASCADE,
    reviewer_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (pr_id, reviewer_id)
);

-- +goose Down
DROP TABLE IF EXISTS assigned_pr_reviewers;

DROP TABLE IF EXISTS pullrequests;

DROP TYPE IF EXISTS pr_status;