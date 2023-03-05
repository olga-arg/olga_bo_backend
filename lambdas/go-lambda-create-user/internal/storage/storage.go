package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go-lambda-create-user/pkg/domain"
)

type UserRepository struct {
	db *dynamodb.DynamoDB
}

func NewUserRepository(db *dynamodb.DynamoDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(user *domain.User) error {
	fmt.Println("Saving user", user)
	fmt.Println("Saving user ID", user.ID)
	fmt.Println("Saving user Name", user.Name)
	fmt.Println("Saving user Email", user.Email)
	item, err := dynamodbattribute.MarshalMap(user)
	fmt.Println("Saving user item", item)
	if err != nil {
		fmt.Println("Error marshalling user", err)
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("usersTable"),
		Item:      item,
	})
	fmt.Println("Error saving user", err)

	return err
}
