package processor

import (
	"bytes"
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/badoux/checkmail"
	"go-lambda-bulk-create-users/pkg/dto"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strings"
)

type Processor interface {
	CreateMultipleUsers(ctx context.Context, inputs []dto.CreateUserInput, companyId string) ([]domain.UserNotCreated, error)
	ValidateUserInput(ctx context.Context, input *dto.CreateUserInput) error
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
	ParseCSVFromRequest(ctx context.Context, request events.APIGatewayProxyRequest) ([]dto.CreateUserInput, error)
}

type processor struct {
	userStorage db.UserRepository
}

func New(s db.UserRepository) Processor {
	return &processor{
		userStorage: s,
	}
}

func (p *processor) CreateMultipleUsers(ctx context.Context, inputs []dto.CreateUserInput, companyId string) ([]domain.UserNotCreated, error) {
	var usersToSave []*domain.User          // Para almacenar usuarios que serán guardados en la base de datos
	var failedUsers []domain.UserNotCreated // Para almacenar usuarios que fallaron en ser creados

	// Configuración inicial de Cognito, como antes
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return failedUsers, err
	}
	cognitoClient := cognitoidentityprovider.New(sess)
	userPoolID := os.Getenv("USER_POOL_ID")

	for _, input := range inputs {
		user, err := domain.NewUser(input.Name, input.Surname, input.Email)
		if err != nil {
			failedUsers = append(failedUsers, domain.UserNotCreated{Email: input.Email, Reason: err.Error()})
			continue
		}
		user.MonthlyLimit = 10 // Ajustar según tu lógica
		role, err := domain.ParseUserRole(input.Role)
		if err != nil {
			failedUsers = append(failedUsers, domain.UserNotCreated{Email: input.Email, Reason: err.Error()})
			continue
		}
		user.Role = role

		createUserInput := &cognitoidentityprovider.AdminCreateUserInput{
			MessageAction: aws.String("SUPPRESS"),
			Username:      aws.String(input.Email),
			UserPoolId:    aws.String(userPoolID),
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{Name: aws.String("email"), Value: aws.String(input.Email)},
				{Name: aws.String("name"), Value: aws.String(companyId)},
				{Name: aws.String("email_verified"), Value: aws.String("False")},
			},
		}

		_, err = cognitoClient.AdminCreateUser(createUserInput)
		if err != nil {
			failedUsers = append(failedUsers, domain.UserNotCreated{Email: input.Email, Reason: err.Error()})
			continue
		}

		usersToSave = append(usersToSave, user)
	}

	if len(usersToSave) > 0 {
		if err := p.userStorage.SaveMultipleUsers(usersToSave, companyId); err != nil {
			// Considera manejar los errores de la base de datos de manera que puedas especificar cuáles usuarios fallaron aquí también
			fmt.Println("Error saving users to database: ", err)
			// Este error no se agrega a failedUsers porque es un fallo en el batch, no individual
			return failedUsers, err
		}
	}

	return failedUsers, nil // Devuelve la lista de usuarios que no se pudieron crear junto con cualquier error global
}

func (p *processor) ValidateUserInput(ctx context.Context, input *dto.CreateUserInput) error {
	fmt.Println("Validating input")
	if input.Name == "" || len(input.Name) < 2 || len(input.Name) > 50 {
		return fmt.Errorf("name validation error")
	}
	if input.Surname == "" || len(input.Surname) < 2 || len(input.Surname) > 50 {
		return fmt.Errorf("surname validation error")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	err := checkmail.ValidateFormat(input.Email)
	if err != nil {
		return fmt.Errorf("invalid email format: %v", err)
	}
	return nil
}

func (p *processor) ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error) {
	// Validate user
	isAuthorized, err := p.userStorage.IsUserAuthorized(email, companyId, allowedRoles)
	if err != nil {
		return false, err
	}
	if isAuthorized {
		return true, nil
	}
	return false, nil
}

func (p *processor) ParseCSVFromRequest(ctx context.Context, request events.APIGatewayProxyRequest) ([]dto.CreateUserInput, error) {
	var reader io.Reader
	if request.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return nil, fmt.Errorf("error decoding base64 body: %v", err)
		}
		reader = bytes.NewReader(decodedBody)
	} else {
		reader = strings.NewReader(request.Body)
	}

	fmt.Println("Content-Type header: ", request.Headers["content-type"])

	contentType := request.Headers["content-type"]
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("error parsing Content-Type header: %v", err)
	}

	fmt.Println("Finished parsing Content-Type header")

	mr := multipart.NewReader(reader, params["boundary"])

	fmt.Println("Created multipart reader")

	var users []dto.CreateUserInput

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next part of multipart request: %v", err)
		}

		if part.FormName() == "csvfile" {
			csvReader := csv.NewReader(part)
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("error reading CSV record: %v", err)
				}

				if len(record) < 4 {
					continue // O maneja el error como prefieras
				}

				user := dto.CreateUserInput{
					Name:    record[0],
					Surname: record[1],
					Email:   record[2],
					Role:    record[3],
				}

				users = append(users, user)
			}
			// Asegúrate de cerrar la parte después de leerla
			part.Close()
			break // Rompe el ciclo después de encontrar y leer el CSV
		}
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no CSV file part found in the request or empty CSV")
	}

	fmt.Println("Users parsed from CSV: ", users)

	return users, nil
}
