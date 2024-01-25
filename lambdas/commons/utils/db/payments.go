package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
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
func getUserTable(companyID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_users", companyID)
		return tx.Table(tableName)
	}
}

func (r *PaymentRepository) GetUserIdByEmail(email string, companyID string) (*domain.User, error) {
	var user domain.User
	err := r.Db.Scopes(getUserTable(companyID)).Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("Error getting user id: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *PaymentRepository) Save(payment *domain.Payment, companyId string) error {
	err := r.Db.Scopes(getPaymentTable(companyId)).Create(payment).Error
	if err != nil {
		fmt.Println("Error saving Payment: ", err)
		return err
	}
	return nil
}

func (r *PaymentRepository) UpdateUser(newUser *domain.User, companyId string) error {
	// Configura GORM para usar la tabla específica basada en companyId
	err := r.Db.Scopes(getUserTable(companyId)).Save(newUser).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
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
