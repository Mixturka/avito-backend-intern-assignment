// В данном пакете объявляются все интерфейсы конкретно сервисов
package service

import (
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
)

type Team interface {
	CreateTeam(ctx context.Context, team entity.Team) error
	GetTeamWithMembers(ctx context.Context, teamName string) (entity.Team, error)
}

type User interface {
	UpdateReposDB(db db.TransactionalDB) User
	SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error)
}

type PullRequest interface {
	Create(ctx context.Context, pr entity.PullRequest) (*entity.PullRequest, error)
	MarkMerged(ctx context.Context, prID string) (*entity.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (*entity.PullRequest, string, error)
	GetPRsByReviewer(ctx context.Context, userID string) (string, []entity.PullRequest, error)
}
