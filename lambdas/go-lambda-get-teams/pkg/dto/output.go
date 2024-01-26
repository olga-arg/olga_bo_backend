package dto

import (
	"commons/domain"
)

type Output struct {
	Teams domain.Teams `json:"teams"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(teams domain.DbTeams) *Output {
	teamMap := make(map[string]*domain.Team)
	var dtoTeams domain.Teams
	//teamsAlreadyVisited := make(map[string]bool)
	for _, team := range teams {
		if existingTeam, ok := teamMap[team.ID]; ok {
			// El equipo ya existe, agregamos el usuario al equipo existente
			user := domain.User{
				ID:   team.UserId,
				Name: team.UserName,
			}
			existingTeam.Users = append(existingTeam.Users, user)
		} else {
			// El equipo no existe, creamos uno nuevo y lo agregamos al mapa
			team := domain.Team{
				ID:   team.ID,
				Name: team.Name,
				// Agregar otros campos de equipo según sea necesario
				Users: domain.Users{
					domain.User{
						ID:              team.UserId,
						Name:            team.UserName,
						Surname:         team.UserSurname,
						Email:           team.UserEmail,
						MonthlySpending: team.UserMonthlySpending,
						// Agregar otros campos de usuario según sea necesario
					},
				},
				ReviewerId:      team.ReviewerId,
				MonthlySpending: team.MonthlySpending,
				AnnualBudget:    team.AnnualBudget,
				Status:          team.Status,
				CreatedDate:     team.CreatedDate,
			}
			dtoTeams = append(dtoTeams, team)
			teamMap[team.ID] = &dtoTeams[len(dtoTeams)-1]
		}
	}
	return &Output{
		Teams: dtoTeams,
	}
}
