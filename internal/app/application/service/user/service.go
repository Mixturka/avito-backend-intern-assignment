package user

import (
	"avito-backend-intern-assignment/internal/app/application/service"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	db.TransactionalRepository[Repository]
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, userID string) (*entity.User, error)
	GetByTeam(ctx context.Context, teamName string) ([]entity.User, error)
}

type Service struct {
	usersRepo Repository
}

func NewService(usersRepo Repository) *Service {
	return &Service{
		usersRepo: usersRepo,
	}
}

func (s *Service) SetIsActive(ctx context.Context, userID string, isActive bool) (entity.User, error) {
	u, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("get user by id: %w", err)
	}

	if u == nil {
		return entity.User{}, ErrUserNotFound
	}

	u.IsActive = isActive
	err = s.usersRepo.Update(ctx, *u)
	if err != nil {
		return entity.User{}, err
	}

	return *u, nil
}

func (s *Service) UpdateReposDB(db db.TransactionalDB) service.User {
	return &Service{
		usersRepo: s.usersRepo.WithDB(db),
	}
}
