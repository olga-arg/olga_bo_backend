package dto

import (
	"commons/domain"
)

type Output struct {
	Teams []domain.TeamOutput `json:"teams"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(teams domain.DbTeams) *Output {
	teamMap := make(map[string]*domain.TeamOutput)
	var dtoTeams []domain.TeamOutput

	for _, dbTeam := range teams {
		var reviewer *domain.User
		if dbTeam.ReviewerId != "" {
			reviewer = &domain.User{
				ID:       dbTeam.ReviewerId,
				Name:     dbTeam.Reviewer.Name,
				Surname:  dbTeam.Reviewer.Surname,
				FullName: dbTeam.Reviewer.FullName,
				Email:    dbTeam.Reviewer.Email,
				Status:   dbTeam.Reviewer.Status,
				Role:     dbTeam.Reviewer.Role,
			}
		} // reviewer es nil si dbTeam.ReviewerId es ""

		if existingTeam, ok := teamMap[dbTeam.ID]; ok {
			if dbTeam.UserId != "" {
				user := domain.User{
					ID:       dbTeam.UserId,
					Name:     dbTeam.UserName,
					Surname:  dbTeam.UserSurname,
					FullName: dbTeam.UserFullName,
					Email:    dbTeam.UserEmail,
					Status:   dbTeam.UserStatus,
					Role:     dbTeam.UserRole,
				}
				existingTeam.Users = append(existingTeam.Users, user)
			}
		} else {
			newTeam := domain.TeamOutput{
				ID:              dbTeam.ID,
				Name:            dbTeam.Name,
				MonthlySpending: dbTeam.MonthlySpending,
				AnnualBudget:    dbTeam.AnnualBudget,
				Status:          dbTeam.Status,
				CreatedDate:     dbTeam.CreatedDate,
				Reviewer:        reviewer,
				ReviewerId:      dbTeam.ReviewerId,
			}

			if dbTeam.UserId != "" {
				newTeam.Users = []domain.User{
					{
						ID:       dbTeam.UserId,
						Name:     dbTeam.UserName,
						Surname:  dbTeam.UserSurname,
						FullName: dbTeam.UserFullName,
						Email:    dbTeam.UserEmail,
						Status:   dbTeam.UserStatus,
						Role:     dbTeam.UserRole,
					},
				}
			} else {
				newTeam.Users = []domain.User{}
			}

			dtoTeams = append(dtoTeams, newTeam)
			teamMap[dbTeam.ID] = &dtoTeams[len(dtoTeams)-1]
		}
	}

	return &Output{
		Teams: dtoTeams,
	}
}
