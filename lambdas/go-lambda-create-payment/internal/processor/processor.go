package processor

import (
	"context"
	"errors"
	"fmt"
	"go-lambda-create-payment/internal/storage"
	"go-lambda-create-payment/pkg/domain"
	"go-lambda-create-payment/pkg/dto"
)

type Processor interface {
	CreatePayment(ctx context.Context, input *dto.CreatePaymentInput) error
}

type processor struct {
	storage     storage.PaymentRepository
	teamStorage storage.TeamRepository
}

func New(paymentRepo storage.PaymentRepository, teamRepo storage.TeamRepository) Processor {
	return &processor{
		storage:     paymentRepo,
		teamStorage: teamRepo,
	}
}

func (p *processor) CreatePayment(ctx context.Context, input *dto.CreatePaymentInput) error {
	var shopName string
	shopName = "Amazon"
	var cardID string
	cardID = "1234567890"
	var userID string
	userID = "1234567890"
	var category string
	category = "Groceries"
	var receipt string
	receipt = "https://s3.amazonaws.com/your-bucket-name/receipt.jpg"
	var paymentType domain.PaymentType
	paymentType = domain.Card

	if input.Amount < 0 {
		fmt.Println("Error: amount cannot be less than 0")
		return errors.New("amount cannot be less than 0")
	}

	// If there is an input.Label then it is a payment to a team
	if input.Label != "" {
		// Find team by ID
		team, err := p.teamStorage.FindTeamByID(input.Label)
		if err != nil {
			fmt.Println("Error finding team: ", err)
			return err
		}
		// If the team exists then set the payment type to team
		err = p.teamStorage.UpdateTeamMonthlySpending(team, input.Amount)
		if err != nil {
			fmt.Println("Error updating payment type to team: ", err)
			return err
		}
	}

	// Creates a new payment
	payment, err := domain.NewPayment(input.Amount, shopName, cardID, userID, category, receipt, input.Label, paymentType)
	if err != nil {
		fmt.Println("Error creating payment: ", err)
		return err
	}
	// Saves the payment to the database if it doesn't already exist
	if err := p.storage.Save(payment); err != nil {
		fmt.Println("Error saving payment: ", err)
		return err
	}
	// Returns
	return nil
}
