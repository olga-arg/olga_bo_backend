package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type PaymentRepository struct {
	Db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	db.AutoMigrate(&domain.Payment{})
	return &PaymentRepository{
		Db: db,
	}
}

func getPaymentTable(companyID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_payments", companyID)
		return tx.Table(tableName)
	}
}

func (r *PaymentRepository) Save(payment *domain.Payment, companyId string) error {
	err := r.Db.Scopes(getPaymentTable(companyId)).Create(payment).Error
	if err != nil {
		fmt.Println("Error saving Payment: ", err)
		return err
	}
	return nil
}

func (r *PaymentRepository) GetAllPayments(filters map[string]string, companyId string) ([]domain.Payment, error) {
	var payments []domain.Payment

	// Start building the query using GORM and the predefined scopes for dynamic table names
	query := r.Db.Scopes(getPaymentTable(companyId)).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(getUserTable(companyId))
		})

	// Apply filters to the query
	if paymentType, ok := filters["payment_type"]; ok {
		query = query.Where("type = ?", paymentType)
	}
	if hasReceipt, ok := filters["receipt"]; ok {
		if hasReceipt == "true" {
			query = query.Where("receipt_image_key <> ''")
		} else if hasReceipt == "false" {
			query = query.Where("receipt_image_key = '' OR receipt_image_key IS NULL")
		}
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	// Order and execute the query
	query = query.Order("created_date")
	err := query.Find(&payments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No payments found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting payments:", err)
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) UpdatePayment(newPayment *domain.Payment, companyId string) error {
	query := r.Db.Scopes(getPaymentTable(companyId)).Save(newPayment)
	if query.Error != nil {
		fmt.Println("Error updating payment:", query.Error)
		return errors.Wrap(query.Error, "failed to update payment")
	}
	return nil
}

func (r *PaymentRepository) GetPaymentByID(paymentID string, companyId string) (*domain.Payment, error) {
	var payment domain.Payment
	query := r.Db.Scopes(getPaymentTable(companyId)).Where("id = ?", paymentID)
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
