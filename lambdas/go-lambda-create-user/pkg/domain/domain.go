package domain

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

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
}

func NewUser(name, surname, email string) (*User, error) {
	var user User
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println("error generating uuid: ", err)
		return nil, err
	}
	user.ID = id.String()
	user.Name = name
	user.Surname = surname
	user.FullName = name + " " + surname
	user.Email = email
	user.Status = Pending
	user.CreatedDate = time.Now()
	return &user, nil
}
