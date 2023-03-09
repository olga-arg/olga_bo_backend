package domain

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

type User struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Surname string             `json:"surname"`
	Email   string             `json:"email"`
	Limit   int                `json:"limit" default:"0"`
	IsAdmin bool               `json:"isAdmin" default:"false"`
	Teams   []string           `json:"team" default:"[]"`
	Status  ConfirmationStatus `json:"status" default:"Pending"`
}

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
