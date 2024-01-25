package domain

import (
	"fmt"
	"github.com/google/uuid"

	"time"
)

type PaymentType int

type Payment struct {
	ID              string             `json:"id"`
	Amount          float32            `json:"amount"`
	ShopName        string             `json:"shop_name"`
	Cuit            string             `json:"cuit"`
	Date            string             `json:"date"`
	Time            string             `json:"time"`
	UserID          string             `json:"user_id"`
	Category        string             `json:"category"`
	receiptNumber   string             `json:"receiptNumber"`
	receiptType     string             `json:"receiptType"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	ReceiptImageKey string             `json:"receiptImageKey"`
	CreatedDate     time.Time          `json:"created"`
}

func NewPayment(amount float32, shopName, cuit, date, _time, category, receiptNumber, receiptType, receiptImageKey, userId string) (*Payment, error) {
	var payment Payment
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	payment.ID = id.String()
	payment.Amount = amount
	payment.ShopName = shopName
	payment.Cuit = cuit
	payment.Date = date
	payment.Time = _time
	payment.UserID = userId
	payment.Category = category
	payment.receiptNumber = receiptNumber
	payment.receiptType = receiptType
	payment.Status = Pending
	payment.ReceiptImageKey = receiptImageKey
	payment.CreatedDate = time.Now()
	return &payment, nil
}
