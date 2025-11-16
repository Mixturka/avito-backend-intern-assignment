package entity

type User struct {
	IsActive bool
	TeamName string
	UserId   string
	Username string
}

func (u User) ToDomainTeamMember() TeamMember {
	return TeamMember{
		IsActive: u.IsActive,
		UserId:   u.UserId,
		Username: u.Username,
	}
}
