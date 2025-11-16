package team

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/application/mappers"
	"avito-backend-intern-assignment/internal/app/application/service"
	"avito-backend-intern-assignment/internal/app/application/service/team"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"context"
	"time"
)

type Handler struct {
	teamService service.Team
}

func NewHandler(teamService service.Team) *Handler {
	return &Handler{
		teamService: teamService,
	}
}

func (h *Handler) GetTeamGet(ctx context.Context, request api.GetTeamGetRequestObject) (api.GetTeamGetResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	var t entity.Team
	var err error
	if t, err = h.teamService.GetTeamWithMembers(serviceCtx, request.Params.TeamName); err != nil {
		if err == team.ErrTeamNotFound {
			return api.GetTeamGet404JSONResponse{}, nil
		}

		return api.GetTeamGet500JSONResponse{}, api.ErrInternalServer
	}

	return api.GetTeamGet200JSONResponse(mappers.ToApiTeam(t)), nil
}

func (h *Handler) PostTeamAdd(ctx context.Context, request api.PostTeamAddRequestObject) (api.PostTeamAddResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	teamEntity := mappers.ToEntityTeam(*request.Body)
	if err := h.teamService.CreateTeam(serviceCtx, teamEntity); err != nil {
		return api.PostTeamAdd500JSONResponse{}, api.ErrInternalServer
	}

	resp := api.PostTeamAdd201JSONResponse{
		Team: request.Body,
	}
	return resp, nil
}
