package domain

import (
    "github.com/google/uuid"
)

type User struct {
    ID    string
    Name  string
    Email string
}

func NewUser(name, email string) *User {
    return &User{
        ID:    uuid.New().String(),
        Name:  name,
        Email: email,
    }
}
