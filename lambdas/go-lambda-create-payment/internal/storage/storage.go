package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-lambda-create-payment/pkg/domain"
)

type PaymentRepository struct {
	Db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		Db: db,
	}
}

func getPaymentTable(payment *domain.Payment) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Extract company name to specify the table name
		return tx.Table("payments")
	}
}

func (r *PaymentRepository) Save(payment *domain.Payment) error {
	err := r.Db.Scopes(getPaymentTable(payment)).AutoMigrate(&domain.Payment{}).Create(payment).Error
	if err != nil {
		fmt.Println("Error saving Payment: ", err)
		return err
	}
	return nil
}
