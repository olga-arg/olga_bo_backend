package dto

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go-lambda-get-all-users/pkg/domain"
	"strconv"
)

type User struct {
	ID      string                    `json:"id"`
	Name    string                    `json:"name"`
	Surname string                    `json:"surname"`
	Email   string                    `json:"email"`
	Limit   int                       `json:"limit"`
	IsAdmin bool                      `json:"isAdmin"`
	Teams   []string                  `json:"team"`
	Status  domain.ConfirmationStatus `json:"status"`
}

type Output struct {
	Users []User `json:"users"`
}

// Unmarshal the DynamoDB response items into User objects
func UnmarshalUsers(items []*dynamodb.AttributeValue) ([]User, error) {
	users := make([]User, 0, len(items))
	for _, item := range items {
		user := User{}
		if id, ok := item.M["id"]; ok {
			user.ID = aws.StringValue(id.S)
		}
		if name, ok := item.M["name"]; ok {
			user.Name = aws.StringValue(name.S)
		}
		if surname, ok := item.M["surname"]; ok {
			user.Surname = aws.StringValue(surname.S)
		}
		if email, ok := item.M["email"]; ok {
			user.Email = aws.StringValue(email.S)
		}
		if limit, ok := item.M["limit"]; ok {
			i, err := strconv.Atoi(aws.StringValue(limit.N))
			if err != nil {
				return nil, err
			}
			user.Limit = i
		}
		if isAdmin, ok := item.M["isAdmin"]; ok {
			user.IsAdmin = aws.BoolValue(isAdmin.BOOL)
		}
		if teams, ok := item.M["team"]; ok {
			teamsList := teams.SS
			user.Teams = make([]string, 0, len(teamsList))
			for _, team := range teamsList {
				user.Teams = append(user.Teams, aws.StringValue(team))
			}
		}
		if status, ok := item.M["status"]; ok {
			user.Status = domain.ParseConfirmationStatus(aws.StringValue(status.S))
		}
		users = append(users, user)
	}
	return users, nil
}
