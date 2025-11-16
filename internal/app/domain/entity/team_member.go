package entity

type TeamMember struct {
	IsActive bool
	UserId   string
	Username string
}

func (tm TeamMember) ToDomainUser(teamName string) User {
	return User{
		IsActive: tm.IsActive,
		UserId:   tm.UserId,
		Username: tm.Username,
		TeamName: teamName,
	}
}
