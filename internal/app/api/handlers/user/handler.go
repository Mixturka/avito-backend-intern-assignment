package user

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/application/mappers"
	"avito-backend-intern-assignment/internal/app/application/service"
	"avito-backend-intern-assignment/internal/app/application/service/user"
	"context"
	"errors"
	"time"
)

type Handler struct {
	userService service.User
}

func NewHandler(userService service.User) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	u, err := h.userService.SetIsActive(serviceCtx, request.Body.UserId, request.Body.IsActive)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return api.PostUsersSetIsActive404JSONResponse{}, nil
		}

		return api.PostUsersSetIsActive500JSONResponse{}, nil
	}

	userDTO := mappers.ToApiUser(u)
	return api.PostUsersSetIsActive200JSONResponse{
		User: &userDTO,
	}, nil
}
