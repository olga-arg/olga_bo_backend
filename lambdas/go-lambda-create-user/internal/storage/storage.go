package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go-lambda-create-user/pkg/domain"
	"log"
)

type UserRepository struct {
	db *dynamodb.DynamoDB
}

func NewUserRepository(db *dynamodb.DynamoDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) EmailAlreadyExists(email string) (bool, error) {
	log.Println("Checking if email already exists: ", email)
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("usersTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		log.Println("Error getting user: ", err)
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) Save(user *domain.User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Println("Error marshalling user", err)
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("usersTable"),
		Item:      item,
	})
	if err != nil {
		log.Println("Error saving user: ", err)
	}
	return err
}
