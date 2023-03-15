package domain

import (
	"fmt"
	"github.com/google/uuid"

	"time"
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Canceled
	Approved
)

type PaymentType int

const (
	Card = iota
	Cash
)

type Payment struct {
	ID          string             `json:"id"`
	Amount      float32            `json:"amount"`
	ShopName    string             `json:"shop_name"`
	CardId      string             `json:"card_id"`
	Type        PaymentType        `json:"payment_type"`
	UserID      string             `json:"user_id"`
	Category    string             `json:"category"`
	Label       string             `json:"label"`
	Status      ConfirmationStatus `json:"status" default:"Pending"`
	Receipt     string             `json:"receipt"`
	CreatedDate time.Time          `json:"created"`
}

func NewPayment(amount float32, shopName, cardID, userID, category, receipt, label string, paymentType PaymentType) (*Payment, error) {
	var payment Payment
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	payment.ID = id.String()
	payment.Amount = amount
	payment.ShopName = shopName
	payment.CardId = cardID
	payment.UserID = userID
	payment.Category = category
	payment.Label = label
	payment.Receipt = receipt
	payment.Type = paymentType
	payment.Status = Pending
	payment.CreatedDate = time.Now()
	return &payment, nil
}
