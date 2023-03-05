package storage

import (
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
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("usersTable"),
		Item:      item,
	})

	return err
}
