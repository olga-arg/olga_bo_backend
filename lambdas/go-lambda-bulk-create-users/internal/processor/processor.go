package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
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
	CreateMultipleUsers(ctx context.Context, inputs []dto.CreateUserInput, companyId string) error
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

func (p *processor) CreateMultipleUsers(ctx context.Context, inputs []dto.CreateUserInput, companyId string) error {
	var usersToSave []*domain.User // Para almacenar usuarios que serán guardados en la base de datos

	// Sesión de AWS y cliente de Cognito (asumiendo que esto no cambia entre llamadas)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}
	cognitoClient := cognitoidentityprovider.New(sess)
	userPoolID := os.Getenv("USER_POOL_ID")

	for _, input := range inputs {
		// Creación de la estructura de usuario (domain.NewUser)
		user, err := domain.NewUser(input.Name, input.Surname, input.Email)
		if err != nil {
			fmt.Println("Error creating user structure: ", err)
			continue // Decide cómo manejar este error; podría ser agregando a un log, etc.
		}
		user.MonthlyLimit = 10 // TODO: Adjust according to your logic

		role, err := domain.ParseUserRole(input.Role)
		if err != nil {
			fmt.Println("Error parsing user role: ", err)
			continue
		}
		user.Role = role

		// Creación del usuario en Cognito
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
			fmt.Println("Error creating user in Cognito: ", err)
			continue
		}

		usersToSave = append(usersToSave, user)
	}

	// Ahora, guarda todos los usuarios en la base de datos en una sola operación
	if len(usersToSave) > 0 {
		if err := p.userStorage.SaveMultipleUsers(usersToSave, companyId); err != nil {
			fmt.Println("Error saving users to database: ", err)
			return err
		}
	}
	return nil
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
	reader := strings.NewReader(request.Body)

	// Crea un multipart reader.
	// Necesitarás el "boundary" para crear el reader, que se encuentra en el header "Content-Type" del request.
	contentType := request.Headers["Content-Type"]
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("error parsing Content-Type header: %v", err)
	}

	mr := multipart.NewReader(reader, params["boundary"])

	// Encuentra la parte del formulario que contiene el archivo CSV.
	var csvPart *multipart.Part
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next part of multipart request: %v", err)
		}

		if part.FormName() == "csvfile" { // Asume que el campo del formulario se llama "csvfile"
			csvPart = part
			break
		}
	}

	if csvPart == nil {
		return nil, fmt.Errorf("no CSV file part found in the request")
	}

	// Ahora que tienes la parte que contiene el CSV, puedes leerla y parsearla.
	r := csv.NewReader(csvPart)

	// Leer la primera fila (encabezados de columna) y descartarla si es necesario
	_, err = r.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	var users []dto.CreateUserInput

	// Continuar leyendo el resto del archivo CSV como antes
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %v", err)
		}

		// Procesamiento de cada fila asumiendo formato: Name, Surname, Email, Role
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

	return users, nil
}
