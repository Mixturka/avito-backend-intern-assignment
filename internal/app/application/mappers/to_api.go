package mappers

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/domain/entity"
)

func ToApiUser(u entity.User) api.User {
	return api.User{
		IsActive: u.IsActive,
		TeamName: u.TeamName,
		UserId:   u.UserId,
		Username: u.Username,
	}
}

func ToApiTeamMember(tm entity.TeamMember) api.TeamMember {
	return api.TeamMember{
		IsActive: tm.IsActive,
		UserId:   tm.UserId,
		Username: tm.Username,
	}
}

func toApiTeamMembers(members []entity.TeamMember) []api.TeamMember {
	result := make([]api.TeamMember, len(members))
	for i, m := range members {
		result[i] = ToApiTeamMember(m)
	}

	return result
}

func ToApiTeam(t entity.Team) api.Team {
	return api.Team{
		Members:  toApiTeamMembers(t.Members),
		TeamName: t.TeamName,
	}
}

func ToApiPullRequest(pr entity.PullRequest) api.PullRequest {
	return api.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            api.PullRequestStatus(pr.Status),
	}
}

func ToApiPullRequests(prs []entity.PullRequest) []api.PullRequest {
	result := make([]api.PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = ToApiPullRequest(pr)
	}
	return result
}

func ToApiPullRequestShort(pr entity.PullRequest) api.PullRequestShort {
	return api.PullRequestShort{
		AuthorId:        pr.AuthorId,
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		Status:          api.PullRequestShortStatus(pr.Status),
	}
}

func ToApiPullRequestsShort(prs []entity.PullRequest) []api.PullRequestShort {
	result := make([]api.PullRequestShort, len(prs))
	for i, pr := range prs {
		result[i] = ToApiPullRequestShort(pr)
	}
	return result
}
