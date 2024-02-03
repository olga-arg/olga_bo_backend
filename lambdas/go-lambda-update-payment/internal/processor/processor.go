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
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/google/uuid"
	"go-lambda-update-payment/pkg/dto"
	"os"
	"time"
)

type Processor interface {
	UpdatePayment(ctx context.Context, newPayment *domain.Payment, companyId string) error
	GetPayment(ctx context.Context, paymentID, companyId string) (*domain.Payment, error)
	ValidatePaymentInput(ctx context.Context, input *dto.UpdatePaymentInput, request events.APIGatewayProxyRequest, companyId, email string) (*domain.Payment, error)
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	paymentStorage db.PaymentRepository
	teamStorage    db.TeamRepository
	userStorage    db.UserRepository
}

func NewProcessor(paymentRepo db.PaymentRepository, teamRepo db.TeamRepository, userRepo db.UserRepository) Processor {
	return &processor{
		paymentStorage: paymentRepo,
		teamStorage:    teamRepo,
		userStorage:    userRepo,
	}
}

func (p *processor) UpdatePayment(ctx context.Context, newPayment *domain.Payment, companyId string) error {
	err := p.paymentStorage.UpdatePayment(newPayment, companyId)
	if err != nil {
		return err
	}
	return nil
}

func (p *processor) GetPayment(ctx context.Context, paymentID, companyId string) (*domain.Payment, error) {
	user, err := p.paymentStorage.GetPaymentByID(paymentID, companyId)
	if err != nil {
		fmt.Println("Error getting payment by ID", err.Error())
		return nil, err
	}
	return user, nil
}

func (p *processor) ValidatePaymentInput(ctx context.Context, input *dto.UpdatePaymentInput, request events.APIGatewayProxyRequest, companyId, email string) (*domain.Payment, error) {
	fmt.Println("Validating input")
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return nil, fmt.Errorf("invalid request body: %s", err.Error())
	}
	fmt.Println("payment_id: ", request.PathParameters["payment_id"])
	payment, err := p.GetPayment(ctx, request.PathParameters["payment_id"], companyId)
	if err != nil {
		fmt.Println("error getting payment", err.Error())
		return nil, fmt.Errorf("failed to get payment")
	}

	// If the payment status is Exported, it cannot be updated
	if payment.Status == domain.Exported {
		return nil, fmt.Errorf("payment is already exported")
	}

	user, err := p.userStorage.GetUserIdByEmail(email, companyId)
	if err != nil {
		fmt.Println("Error getting user: ", err)
		return nil, err
	}

	if input.Amount != nil {
		payment.Amount = *input.Amount

		// Update the monthly spending of the user
		userTeams, err := p.teamStorage.GetTeamByUserID(user.ID, companyId)
		fmt.Println("User teams: ", userTeams)
		if err != nil {
			fmt.Println("Error getting user teams ", err)
			return nil, err
		}
		// If the user is part of a team, update the team's monthly spending
		if len(userTeams) > 0 {
			fmt.Println("Updating team monthly spending")
			for _, userTeam := range userTeams {
				team, err := p.teamStorage.GetTeamByID(userTeam.TeamID, companyId)
				if err != nil {
					fmt.Println("Error getting team: ", err)
					return nil, err
				}
				amount := *input.Amount
				team.MonthlySpending += int(amount)
				if err := p.teamStorage.UpdateTeamMonthlySpending(team.MonthlySpending, companyId); err != nil {
					fmt.Println("Error updating team monthly spending: ", err)
					return nil, err
				}
			}
		}
	}
	// Validate the status has to be: Pending, Approved, Deleted, Exported
	if input.Status != "" {
		var err error
		payment.Status, err = domain.ParseConfirmationStatus(input.Status)
		if err != nil {
			// The status is invalid
			return nil, err
		}

		// Additional check to disallow 'Created' or 'Confirmed' status
		if payment.Status == domain.Created || payment.Status == domain.Confirmed {
			return nil, fmt.Errorf("invalid status: %s", input.Status)
		}
	}

	if input.ShopName != "" {
		payment.ShopName = input.ShopName
	}
	if input.Category != "" {
		payment.Category = input.Category
	}
	if input.Date != (time.Time{}) {
		payment.Date = input.Date
	}

	cuitChanged := input.Cuit != "" && input.Cuit != payment.Cuit
	if cuitChanged {
		payment.Cuit = input.Cuit
		err := logCuitChange(ctx, payment.ID, payment.Cuit, input.Cuit, payment.ReceiptImageKey)
		if err != nil {
			fmt.Println("error logging CUIT change: ", err.Error())
			// Decidir si quieres devolver el error o simplemente registrarlo
		}
	}
	fmt.Println("Input validated successfully")
	return payment, nil
}

func logCuitChange(ctx context.Context, paymentId, oldCuit, newCuit, receiptImageKey string) error {
	stage := os.Getenv("STAGE")
	// Crear una sesión de AWS
	sess := session.Must(session.NewSession())
	// Crear un cliente de CloudWatch Logs
	cw := cloudwatchlogs.New(sess)

	logGroupName := stage + "-cuit-change"
	logStreamName := uuid.New().String()

	// Intentar crear el stream de log y manejar cualquier error
	_, err := cw.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		fmt.Printf("Error al crear el stream de log: %s\n", err)
		return err
	}

	// Crear el mensaje de log
	logMessage := fmt.Sprintf("Cuit changed from %s to %s for payment ID %s. Receipt Image Key: %s", oldCuit, newCuit, paymentId, receiptImageKey)
	timestamp := aws.Int64(time.Now().UnixNano() / int64(time.Millisecond))

	// Preparar los eventos de log
	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(logMessage),
				Timestamp: timestamp,
			},
		},
	}

	// Enviar el evento de log
	_, err = cw.PutLogEvents(input)
	if err != nil {
		return err
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
