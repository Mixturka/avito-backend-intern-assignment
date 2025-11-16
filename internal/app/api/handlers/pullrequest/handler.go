package pullrequest

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/application/mappers"
	"avito-backend-intern-assignment/internal/app/application/service"
	"avito-backend-intern-assignment/internal/app/application/service/pullrequest"
	"avito-backend-intern-assignment/internal/app/application/service/user"
	"context"
	"errors"
	"log"
	"time"
)

type Handler struct {
	prService service.PullRequest
}

func NewHandler(prService service.PullRequest) *Handler {
	return &Handler{
		prService: prService,
	}
}

func (h *Handler) PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	createDTO := api.PostPullRequestCreateJSONBody(*request.Body)
	prEntity := mappers.ToEntityPullRequestCreate(createDTO)

	log.Printf("Handler: Creating PR with ID: %s, Author: %s", prEntity.PullRequestId, prEntity.AuthorId)

	createdPR, err := h.prService.Create(serviceCtx, prEntity)
	if err != nil {
		log.Printf("Handler: PR creation failed: %v", err)

		switch {
		case errors.Is(err, pullrequest.ErrAuthorNotFound):
			return api.PostPullRequestCreate404JSONResponse{}, nil
		case errors.Is(err, pullrequest.ErrPullRequestExists):
			return api.PostPullRequestCreate409JSONResponse{}, nil
		case errors.Is(err, pullrequest.ErrNotEnoughReviewers):
			return api.PostPullRequestCreate409JSONResponse{}, nil
		default:
			log.Printf("Handler: Internal server error during PR creation: %v", err)
			return api.PostPullRequestCreate500JSONResponse{}, nil
		}
	}

	log.Printf("Handler: Created PR - ID: %s, Reviewers: %v",
		createdPR.PullRequestId, createdPR.AssignedReviewers)

	response := mappers.ToApiPullRequest(*createdPR)

	log.Printf("Handler: API response - ID: %s, Reviewers: %v",
		response.PullRequestId, response.AssignedReviewers)

	log.Printf("Handler: Successfully created PR: %s", createdPR.PullRequestId)
	return api.PostPullRequestCreate201JSONResponse{
		Pr: &response,
	}, nil
}

func (h *Handler) PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	prEntity, err := h.prService.MarkMerged(serviceCtx, request.Body.PullRequestId)
	if err != nil {
		if err == pullrequest.ErrPullRequestNotFound {
			return api.PostPullRequestMerge404JSONResponse{}, nil
		}

		return api.PostPullRequestMerge500JSONResponse{}, api.ErrInternalServer
	}

	prDTO := mappers.ToApiPullRequest(*prEntity)
	return api.PostPullRequestMerge200JSONResponse{
		Pr: &prDTO,
	}, nil
}

func (h *Handler) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	prEntity, newAssignedUserID, err := h.prService.ReassignReviewer(serviceCtx, request.Body.PullRequestId, request.Body.OldUserId)
	if err != nil {
		switch err {
		case pullrequest.ErrPullRequestNotFound, user.ErrUserNotFound:
			return api.PostPullRequestReassign404JSONResponse{}, nil
		case pullrequest.ErrReassignViolation:
			return api.PostPullRequestReassign409JSONResponse{}, nil
		default:
			return api.PostPullRequestReassign500JSONResponse{}, api.ErrInternalServer
		}
	}

	prDTO := mappers.ToApiPullRequest(*prEntity)
	return api.PostPullRequestReassign200JSONResponse{
		ReplacedBy: newAssignedUserID,
		Pr:         prDTO,
	}, nil
}

func (h *Handler) GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error) {
	serviceCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	userID, prs, err := h.prService.GetPRsByReviewer(serviceCtx, request.Params.UserId)
	if err != nil {
		return api.GetUsersGetReview500JSONResponse{}, api.ErrInternalServer
	}

	prShortDtos := mappers.ToApiPullRequestsShort(prs)
	return api.GetUsersGetReview200JSONResponse{
		UserId:       userID,
		PullRequests: prShortDtos,
	}, nil
}
