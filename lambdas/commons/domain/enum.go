package domain

const (
	Card = iota
	Cash
)

type ConfirmationStatus int

const (
	Pending   ConfirmationStatus = iota // User, Payment
	Approved                            // Payment
	Created                             // Payment, Team
	Deleted                             // User, Payment, Team
	Exported                            // Payment
	Confirmed                           // User
)

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Approved":
		return Approved
	case "Created":
		return Created
	case "Deleted":
		return Deleted
	case "Exported":
		return Exported
	case "Confirmed":
		return Confirmed
	default:
		return Created
	}
}

type UserRoles int

const (
	Employee UserRoles = iota
	Reviewer
	Admin
	Accountant
)
