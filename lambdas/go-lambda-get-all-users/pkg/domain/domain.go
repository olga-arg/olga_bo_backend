package domain

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

type User struct {
	ID           string             `json:"id"`
	Company      string             `json:"company"`
	Name         string             `json:"name"`
	Surname      string             `json:"surname"`
	Email        string             `json:"email"`
	AccountLimit int                `json:"limit" default:"0"`
	Teams        []string           `json:"teams"`
	IsAdmin      bool               `json:"isAdmin" default:"false"`
	Status       ConfirmationStatus `json:"status" default:"Pending"`
}

type Users []User

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Confirmed":
		return Confirmed
	case "Deleted":
		return Deleted
	default:
		return Pending
	}
}
