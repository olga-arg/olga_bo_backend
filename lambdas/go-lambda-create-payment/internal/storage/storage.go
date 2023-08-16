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
	db.AutoMigrate(&domain.Payment{})
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

func getUserTable(user *domain.User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *PaymentRepository) GetUserIdByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.Db.Scopes(getUserTable(&user)).Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("Error getting user id: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *PaymentRepository) Save(payment *domain.Payment) error {
	err := r.Db.Scopes(getPaymentTable(payment)).Create(payment).Error
	if err != nil {
		fmt.Println("Error saving Payment: ", err)
		return err
	}
	return nil
}

//func (r *TeamRepository) FindTeamByID(id string) (*domain.Team, error) {
//	var team domain.Team
//	err := r.Db.Scopes(getTeamsTable(&team)).Where("id = ?", id).First(&team).Error
//	if err != nil {
//		fmt.Println("Error finding team: ", err)
//		return nil, err
//	}
//	return &team, nil
//}

//func (r *TeamRepository) UpdateTeamMonthlySpending(team *domain.Team, paymentAmount float32) error {
//	// Get the current monthly spending of the team
//	var currentMonthlySpending float32
//	currentMonthlySpending = team.MonthlySpending
//
//	// Add the new payment amount to the current monthly spending
//	var newMonthlySpending float32
//	newMonthlySpending = currentMonthlySpending + paymentAmount
//
//	// Save the new monthly spending to the team
//	err := r.Db.Scopes(getTeamsTable(team)).Model(&team).Update("monthly_spending", newMonthlySpending).Error
//	if err != nil {
//		fmt.Println("Error updating team: ", err)
//		return err
//	}
//	return nil
//}
