package domain

import "fmt"

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

func ParseConfirmationStatus(s string) (ConfirmationStatus, error) {
	switch s {
	case "Pending":
		return Pending, nil
	case "Approved":
		return Approved, nil
	case "Created":
		return Created, nil
	case "Deleted":
		return Deleted, nil
	case "Exported":
		return Exported, nil
	case "Confirmed":
		return Confirmed, nil
	default:
		return -1, fmt.Errorf("invalid status: %s", s)
	}
}

func ConfirmationStatusToString(status ConfirmationStatus) string {
	switch status {
	case Pending:
		return "Pendiente"
	case Approved:
		return "Aprobado"
	case Created:
		return "Creado"
	case Deleted:
		return "Eliminado"
	case Exported:
		return "Exportado"
	case Confirmed:
		return "Confirmado"
	default:
		return "Desconocido"
	}
}

type UserRoles int

const (
	Employee UserRoles = iota
	Reviewer
	Admin
	Accountant
)
