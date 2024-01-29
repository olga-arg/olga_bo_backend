package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-payment/pkg/dto"
	"time"
)

type Processor interface {
	UpdatePayment(ctx context.Context, newPayment *domain.Payment, companyId string) error
	GetPayment(ctx context.Context, paymentID, companyId string) (*domain.Payment, error)
	ValidatePaymentInput(ctx context.Context, input *dto.UpdatePaymentInput, request events.APIGatewayProxyRequest, companyId string) (*domain.Payment, error)
}

type processor struct {
	paymentStorage *db.PaymentRepository
}

func NewProcessor(storage *db.PaymentRepository) Processor {
	return &processor{
		paymentStorage: storage,
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

func (p *processor) ValidatePaymentInput(ctx context.Context, input *dto.UpdatePaymentInput, request events.APIGatewayProxyRequest, companyId string) (*domain.Payment, error) {
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
	if input.Amount != nil {
		payment.Amount = *input.Amount
	}
	if input.Status != nil {
		payment.Status = domain.ConfirmationStatus(*input.Status)
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
	fmt.Println("Input validated successfully")
	return payment, nil
}
