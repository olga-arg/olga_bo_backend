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
