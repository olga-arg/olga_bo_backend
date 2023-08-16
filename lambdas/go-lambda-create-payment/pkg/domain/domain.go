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

type Team struct {
	ID              string             `json:"id"`
	CompanyID       string             `json:"company"`
	Name            string             `json:"name"`
	Users           []*User            `gorm:"many2many:user_teams;"`
	ReviewerId      string             `json:"reviewer_id"`
	AnnualBudget    int                `json:"annual_budget"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

type User struct {
	ID              string             `json:"id"`
	CompanyID       string             `json:"company"`
	Name            string             `json:"name"`
	Surname         string             `json:"surname"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	PurchaseLimit   int                `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                `json:"monthly_limit" default:"0"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	IsAdmin         bool               `json:"isAdmin" default:"false"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
	Teams           []*Team            `gorm:"many2many:user_teams;"`
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
