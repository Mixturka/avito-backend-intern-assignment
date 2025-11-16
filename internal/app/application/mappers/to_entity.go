package mappers

import (
	"avito-backend-intern-assignment/internal/app/api"
	"avito-backend-intern-assignment/internal/app/domain/entity"
)

func ToEntityUser(u api.User) entity.User {
	return entity.User{
		IsActive: u.IsActive,
		TeamName: u.TeamName,
		UserId:   u.UserId,
		Username: u.Username,
	}
}

func ToEntityTeamMember(tm api.TeamMember) entity.TeamMember {
	return entity.TeamMember{
		IsActive: tm.IsActive,
		UserId:   tm.UserId,
		Username: tm.Username,
	}
}

func ToEntityTeamMembers(members []api.TeamMember) []entity.TeamMember {
	result := make([]entity.TeamMember, len(members))
	for i, m := range members {
		result[i] = ToEntityTeamMember(m)
	}
	return result
}

func ToEntityTeam(t api.Team) entity.Team {
	return entity.Team{
		Members:  ToEntityTeamMembers(t.Members),
		TeamName: t.TeamName,
	}
}

func ToEntityPullRequest(pr api.PullRequest) entity.PullRequest {
	return entity.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            entity.PRStatus(pr.Status),
	}
}

func ToEntityPullRequests(prs []api.PullRequest) []entity.PullRequest {
	result := make([]entity.PullRequest, len(prs))
	for i, pr := range prs {
		result[i] = ToEntityPullRequest(pr)
	}
	return result
}

func ToEntityPullRequestCreate(prReq api.PostPullRequestCreateJSONBody) entity.PullRequest {
	return entity.PullRequest{
		AuthorId:        prReq.AuthorId,
		PullRequestId:   prReq.PullRequestId,
		PullRequestName: prReq.PullRequestName,
	}
}
