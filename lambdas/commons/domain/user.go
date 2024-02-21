package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Surname         string             `json:"surname"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	PurchaseLimit   int                `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                `json:"monthly_limit" default:"0"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
	Teams           []*Team            `gorm:"many2many:user_teams;"`
	Role            UserRoles          `json:"role" default:"Employee"`
}

type Users []User

func NewUser(name, surname, email string) (*User, error) {
	var user User
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	user.ID = id.String()
	//user.CompanyID = ""
	user.Name = name
	user.Surname = surname
	user.FullName = name + " " + surname
	user.Email = email
	user.Status = Pending
	user.CreatedDate = time.Now()
	return &user, nil
}

type UpdateUserRequest struct {
	CompanyId string `json:"company_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type UserNotCreated struct {
	Email  string
	Reason string
}

type CreateUserResult struct {
	FailedUsers  []UserNotCreated
	SuccessCount int
}
