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
	Created
	Deleted
	Exported
)

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Canceled":
		return Canceled
	case "Approved":
		return Approved
	case "Created":
		return Created
	case "Deleted":
		return Deleted
	case "Exported":
		return Exported
	default:
		return Created
	}
}
