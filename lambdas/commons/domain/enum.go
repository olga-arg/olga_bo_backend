package domain

const (
	Card = iota
	Cash
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Canceled
	Approved
)

const (
	Created ConfirmationStatus = iota
	Deleted
	Awating
)

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Created":
		return Created
	case "Deleted":
		return Deleted
	default:
		return Created
	}
}
