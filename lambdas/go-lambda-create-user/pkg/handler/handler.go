package handler

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
)

func CreateUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateUserInput
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	user := domain.User{
		Name:  input.Name,
		Email: input.Email,
	}

	db := application.NewDynamoDBClient()
	userRepository := storage.NewUserRepository(db)

	err = userRepository.Save(&user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating user",
		}, nil
	}

	output := dto.CreateUserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	response, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating user",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(response),
	}, nil
}
