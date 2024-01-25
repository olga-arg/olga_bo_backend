package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Company struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Status     ConfirmationStatus `json:"status" default:"Approved"`
	UsersLimit int                `json:"users_limit" default:"0"`
}

func NewCompany(name string) (*Company, error) {
	var company Company
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	company.ID = id.String()
	//user.CompanyID = ""
	company.Name = name
	return &company, nil
}
