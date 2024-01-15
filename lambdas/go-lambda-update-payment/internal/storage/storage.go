package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-update-payment/pkg/domain"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	db.AutoMigrate(&domain.Payment{})
	return &PaymentRepository{
		db: db,
	}
}

func getPaymentTable(payment *domain.Payment) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Extract company name to specify the table name
		return tx.Table("payments")
	}
}

func getUserTable(paymentID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *PaymentRepository) UpdatePayment(newPayment *domain.Payment) error {
	// Save the updated user
	query := r.db.Save(newPayment)
	if query.Error != nil {
		fmt.Println("Error updating payment:", query.Error)
		return errors.Wrap(query.Error, "failed to update payment")
	}
	return nil
}

func (r *PaymentRepository) GetPaymentByID(paymentID string) (*domain.Payment, error) {
	var payment domain.Payment
	query := r.db.Scopes(getPaymentTable(&payment)).Where("id = ?", paymentID)
	err := query.First(&payment).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrap(err, "payment not found")
		}
		fmt.Println("Error getting payment by ID:", err)
		return nil, errors.Wrap(err, "failed to get payment by ID")
	}
	return &payment, nil
}
