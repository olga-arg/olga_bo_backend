package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"go-lambda-bulk-create-users-process-queue/pkg/dto"
	"os"
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput, companyId string) error
}

type processor struct {
	userStorage db.UserRepository
}

func New(s db.UserRepository) Processor {
	return &processor{
		userStorage: s,
	}
}

func (p *processor) CreateUser(ctx context.Context, input *dto.CreateUserInput, companyId string) error {
	// Creates a new user. New user takes a name and email and returns a user struct
	user, err := domain.NewUser(input.Name, input.Surname, input.Email)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return err
	}

	// TODO: Erase this line when the monthly limit is implemented
	user.MonthlyLimit = 10

	// Parse the role from the input
	role, err := domain.ParseUserRole(input.Role)
	if err != nil {
		return err
	}
	user.Role = role

	// Creates the user in cognito
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	// Create a new CognitoIdentityProvider client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Specify the user pool ID
	userPoolID := os.Getenv("USER_POOL_ID")

	// Create a user in Cognito
	createUserInput := &cognitoidentityprovider.AdminCreateUserInput{
		MessageAction: aws.String("SUPPRESS"),
		Username:      aws.String(input.Email),
		UserPoolId:    aws.String(userPoolID),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Email),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(companyId),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("False"),
			},
		},
	}

	// Call the AdminCreateUser API
	_, err = cognitoClient.AdminCreateUser(createUserInput)
	if err != nil {
		return err
	}
	fmt.Println("User created successfully in Cognito")

	// Saves the user to the database if it doesn't already exist
	if err := p.userStorage.Save(user, companyId); err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	fmt.Println("User created successfully in Cognito and DynamoDB")
	// Returns
	return nil
}
