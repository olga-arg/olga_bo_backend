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
	Date            time.Time          `json:"date"`
	Time            string             `json:"time"`
	UserID          string             `json:"user_id"`
	Category        string             `json:"category"`
	ReceiptNumber   string             `json:"receiptNumber"`
	ReceiptType     string             `json:"receiptType"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	ReceiptImageKey string             `json:"receiptImageKey"`
	CreatedDate     time.Time          `json:"created"`
	User            User               `gorm:"foreignKey:user_id"`
	ImageURL        string             `json:"image_url"`
}

type Payments []Payment

func NewPayment(amount float32, shopName, cuit, _time, category, receiptNumber, receiptType, receiptImageKey, userId string, date time.Time) (*Payment, error) {
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
	payment.ReceiptNumber = receiptNumber
	payment.ReceiptType = receiptType
	payment.Status = Pending
	payment.ReceiptImageKey = receiptImageKey
	payment.CreatedDate = time.Now()
	return &payment, nil
}
