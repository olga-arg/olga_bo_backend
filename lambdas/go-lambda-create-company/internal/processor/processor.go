package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/badoux/checkmail"
	"go-lambda-create-user/pkg/dto"
	"os"
)

type Processor interface {
	CreateCompany(ctx context.Context, input *dto.CreateCompanyInput) error
	ValidateCompanyInput(ctx context.Context, input *dto.CreateCompanyInput, request events.APIGatewayProxyRequest) error
}

type processor struct {
	companyStorage  db.CompanyRepository
	userStorage     db.UserRepository
	categoryStorage db.CategoryRepository
}

func New(c db.CompanyRepository, u db.UserRepository, ca db.CategoryRepository) Processor {
	return &processor{
		userStorage:     u,
		companyStorage:  c,
		categoryStorage: ca,
	}
}

func (p *processor) CreateCompany(ctx context.Context, input *dto.CreateCompanyInput) error {
	// Creates a new company.
	company, err := domain.NewCompany(input.CompanyName)
	if err != nil {
		fmt.Println("Error creating company: ", err)
		return err
	}

	if input.Cuit == "" {
		return fmt.Errorf("cuit is required")
	}

	// Saves the company to the database if it doesn't already exist
	if err := p.companyStorage.Save(company); err != nil {
		fmt.Println("Error saving company: ", err)
		return err
	}

	// Creates users table for the company
	if err := p.companyStorage.CreateCompanySpecificTables(company.ID); err != nil {
		fmt.Println("Error creating users table: ", err)
		return err
	}

	// Create company expense categories
	var categories = map[string]map[string]string{
		"Comidas y Bebidas": {
			"icon":  "mdiSilverwareForkKnife",
			"color": "FF6384",
		},
		"Transporte": {
			"icon":  "mdiSubwayVariant",
			"color": "36A2EB",
		},
		"Electrónica": {
			"icon":  "mdiLaptop",
			"color": "FFCE56",
		},
		"Entretenimiento": {
			"icon":  "mdiRobotHappy",
			"color": "4BC0C0",
		},
		"Material de Oficina": {
			"icon":  "mdiOfficeBuilding",
			"color": "9966FF",
		},
		"Indumentaria": {
			"icon":  "mdiTshirtCrew",
			"color": "FF9F40",
		},
		"Salud y cuidado personal": {
			"icon":  "mdiBottleTonicPlus",
			"color": "C9CBCF",
		},
		"Educación": {
			"icon":  "mdiSchool",
			"color": "7E7F9A",
		},
		"Mascotas": {
			"icon":  "mdiPaw",
			"color": "f5e050",
		},
		"Supermercado": {
			"icon":  "mdiStore",
			"color": "FFC0CB",
		},
		"Viajes": {
			"icon":  "mdiAirplaneTakeoff",
			"color": "1E90FF",
		},
		"Servicios profesionales": {
			"icon":  "mdiAccountWrench",
			"color": "DAA520",
		},
		"Impuestos": {
			"icon":  "mdiCash",
			"color": "B22222",
		},
		"Cuentas y Servicios": {
			"icon":  "mdiAccountCash",
			"color": "FFD700",
		},
		"Donaciones": {
			"icon":  "mdiHandHeart",
			"color": "32CD32",
		},
		"Inversiones": {
			"icon":  "mdiFinance",
			"color": "4682B4",
		},
		"Préstamos y financiación": {
			"icon":  "mdiCurrencyUsd",
			"color": "DA70D6",
		},
		"Suscripciones": {
			"icon":  "mdiCart",
			"color": "40E0D0",
		},
		"Shopping": {
			"icon":  "mdiShopping",
			"color": "FF4500",
		},
		"Otros": {
			"icon":  "mdiMenu",
			"color": "808080",
		},
	}
	for category, data := range categories {
		cat, _ := domain.NewCategory(company.ID, category, data["color"], data["icon"])
		if err := p.categoryStorage.Save(cat); err != nil {
			fmt.Println("Error creating company expense icons: ", err)
			return err
		}
	}

	// Creates a new user. New user takes a name and email and returns a user struct
	user, err := domain.NewUser(input.UserName, input.UserSurname, input.UserEmail)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return err
	}

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
		Username:      aws.String(input.UserEmail),
		UserPoolId:    aws.String(userPoolID),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.UserEmail),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("True"),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(company.ID),
			},
		},
	}

	// Call the AdminCreateUser API
	_, err = cognitoClient.AdminCreateUser(createUserInput)
	if err != nil {
		return err
	}
	fmt.Println("User created successfully in Cognito")

	// TODO: Erase this line when the monthly limit is implemented
	user.MonthlyLimit = 10
	user.Role = domain.Admin

	// Saves the user to the database if it doesn't already exist
	if err := p.userStorage.Save(user, company.ID); err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	fmt.Println("User created successfully in Cognito and DynamoDB")
	// Returns
	return nil
}

func (p *processor) ValidateCompanyInput(ctx context.Context, input *dto.CreateCompanyInput, request events.APIGatewayProxyRequest) error {
	fmt.Println("Validating input")
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.CompanyName == "" {
		return fmt.Errorf("company name is required")
	}
	if input.UserName == "" {
		return fmt.Errorf("name is required")
	}
	if input.UserSurname == "" {
		return fmt.Errorf("surname is required")
	}
	if input.UserEmail == "" {
		return fmt.Errorf("email is required")
	}
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	err := checkmail.ValidateFormat(input.UserEmail)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
