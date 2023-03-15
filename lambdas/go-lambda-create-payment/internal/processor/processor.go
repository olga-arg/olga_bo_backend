package processor

import (
	"context"
	"fmt"
	"go-lambda-create-payment/internal/storage"
	"go-lambda-create-payment/pkg/domain"
)

type Processor interface {
	CreatePayment(ctx context.Context) error
}

type processor struct {
	storage storage.PaymentRepository
}

func New(s storage.PaymentRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreatePayment(ctx context.Context) error {
	var amount float32
	amount = 100.00
	var shopName string
	shopName = "Amazon"
	var cardID string
	cardID = "1234567890"
	var userID string
	userID = "1234567890"
	var category string
	category = "Groceries"
	var label string
	label = "a408bf50-c2cb-11ed-a27e-368de2526bc1"
	var receipt string
	receipt = "https://s3.amazonaws.com/your-bucket-name/receipt.jpg"
	var paymentType domain.PaymentType
	paymentType = domain.Card

	// Creates a new payment
	payment, err := domain.NewPayment(amount, shopName, cardID, userID, category, receipt, label, paymentType)
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
