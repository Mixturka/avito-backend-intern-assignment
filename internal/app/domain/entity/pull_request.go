package entity

import "time"

type PRStatus string

const (
	PullRequestStatusMERGED PRStatus = "MERGED"
	PullRequestStatusOPEN   PRStatus = "OPEN"
)

type PullRequest struct {
	AssignedReviewers []string
	AuthorId          string
	CreatedAt         *time.Time
	MergedAt          *time.Time
	PullRequestId     string
	PullRequestName   string
	Status            PRStatus
}
