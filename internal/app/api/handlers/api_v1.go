package handlers

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/api/handlers/pullrequest"
	"avito-backend-intern-assignment/internal/app/api/handlers/team"
	"avito-backend-intern-assignment/internal/app/api/handlers/user"
	"context"
)

type ApiV1 struct {
	teamHandler *team.Handler
	userHandler *user.Handler
	prHandler   *pullrequest.Handler
}

func NewApiV1(th *team.Handler, uh *user.Handler, prh *pullrequest.Handler) *ApiV1 {
	return &ApiV1{
		teamHandler: th,
		userHandler: uh,
		prHandler:   prh,
	}
}

func (av *ApiV1) GetTeamGet(ctx context.Context, request api.GetTeamGetRequestObject) (api.GetTeamGetResponseObject, error) {
	return av.teamHandler.GetTeamGet(ctx, request)
}

func (av *ApiV1) GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error) {
	return av.prHandler.GetUsersGetReview(ctx, request)
}

func (av *ApiV1) PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error) {
	return av.prHandler.PostPullRequestCreate(ctx, request)
}

func (av *ApiV1) PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error) {
	return av.prHandler.PostPullRequestMerge(ctx, request)
}

func (av *ApiV1) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	return av.prHandler.PostPullRequestReassign(ctx, request)
}

func (av *ApiV1) PostTeamAdd(ctx context.Context, request api.PostTeamAddRequestObject) (api.PostTeamAddResponseObject, error) {
	return av.teamHandler.PostTeamAdd(ctx, request)
}

func (av *ApiV1) PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error) {
	return av.userHandler.PostUsersSetIsActive(ctx, request)
}

var _ api.StrictServerInterface = (*ApiV1)(nil)
