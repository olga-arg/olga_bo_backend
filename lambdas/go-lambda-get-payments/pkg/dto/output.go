package dto

import (
	"go-lambda-get-payments/pkg/domain"
	"time"
)

type Payment struct {
	ID              string                    `json:"id"`
	Amount          float32                   `json:"amount"`
	ShopName        string                    `json:"shop_name"`
	CardId          string                    `json:"card_id"`
	Type            domain.PaymentType        `json:"payment_type"`
	UserID          string                    `json:"user_id"`
	Category        string                    `json:"category"`
	Label           string                    `json:"label"`
	Status          domain.ConfirmationStatus `json:"status" default:"Pending"`
	ReceiptImageKey string                    `json:"receipt_image_key"`
	CreatedDate     time.Time                 `json:"created"`
	User            User                      `gorm:"foreignKey:user_id"`
}

type Output struct {
	Payments []Payment `json:"payments"`
}

type User struct {
	ID              string  `json:"id"`
	CompanyID       string  `json:"company"`
	Name            string  `json:"name"`
	Surname         string  `json:"surname"`
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
	PurchaseLimit   int     `json:"purchase_limit" default:"0"`
	MonthlyLimit    int     `json:"monthly_limit" default:"0"`
	MonthlySpending float32 `json:"monthly_spending" default:"0"`
	IsAdmin         bool    `json:"isAdmin" default:"false"`
}

func NewOutput(payments []domain.Payment) *Output {
	var dtoPayments []Payment
	for _, payment := range payments {
		dtoPayments = append(dtoPayments, Payment{
			ID:              payment.ID,
			Amount:          payment.Amount,
			ShopName:        payment.ShopName,
			Type:            payment.Type,
			Category:        payment.Category,
			Label:           payment.Label,
			Status:          payment.Status,
			ReceiptImageKey: payment.ReceiptImageKey,
			CreatedDate:     payment.CreatedDate,
			User: User{
				CompanyID:       payment.User.CompanyID,
				Name:            payment.User.Name,
				Surname:         payment.User.Surname,
				FullName:        payment.User.FullName,
				Email:           payment.User.Email,
				PurchaseLimit:   payment.User.PurchaseLimit,
				MonthlyLimit:    payment.User.MonthlyLimit,
				MonthlySpending: payment.User.MonthlySpending,
				IsAdmin:         payment.User.IsAdmin,
			},
		})
	}
	return &Output{
		Payments: dtoPayments,
	}
}
