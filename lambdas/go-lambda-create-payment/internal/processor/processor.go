package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
	"go-lambda-create-payment/pkg/dto"
)

type Processor interface {
	CreatePayment(ctx context.Context, input *dto.CreatePaymentInput, email string) error
}

////// Remove storage and import it from commons

type processor struct {
	storage db.PaymentRepository
}

func New(paymentRepo db.PaymentRepository) Processor {
	return &processor{
		storage: paymentRepo,
	}
}

func (p *processor) CreatePayment(ctx context.Context, input *dto.CreatePaymentInput, email string) error {
	// Validate the status of the user (active or not)
	user, err := p.storage.GetUserIdByEmail(email)
	if err != nil {
		fmt.Println("Error getting user: ", err)
		return err
	}

	// Validate the purchase limit
	purchaseLimit := user.PurchaseLimit
	if float32(purchaseLimit) < input.Amount {
		return fmt.Errorf("the amount is greater than the purchase limit")
	}
	// Validate the monthly limit
	monthlyLimit := user.MonthlyLimit
	remainingMonthlyLimit := float32(monthlyLimit) - user.MonthlySpending
	if remainingMonthlyLimit < input.Amount {
		return fmt.Errorf("Error: The amount is greater than the monthly limit")
	}

	// Create payment
	payment, err := domain.NewPayment(input.Amount, input.ShopName, input.Cuit, input.Date, input.Time, input.Category, input.ReceiptNumber, input.ReceiptType, input.ReceiptImageKey, user.ID)
	if err != nil {
		fmt.Println("Error creating payment: ", err)
		return err
	}

	// Update the monthly spending of the user
	user.MonthlySpending += input.Amount
	if err := p.storage.UpdateUser(user); err != nil {
		fmt.Println("Error updating user: ", err)
		return err
	}
	// Save payment to db
	if err := p.storage.Save(payment); err != nil {
		fmt.Println("Error saving payment: ", err)
		return err
	}

	return nil
}
