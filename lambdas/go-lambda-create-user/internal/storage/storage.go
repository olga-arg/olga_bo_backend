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

// The repository is responsible for interacting with the database
func NewUserRepository(db *dynamodb.DynamoDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) EmailAlreadyExists(email string) (bool, error) {
	result, err := r.db.Query(&dynamodb.QueryInput{
		TableName:              aws.String("users"),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		log.Println("Error getting user: ", err)
		return false, err
	}
	// The query should return 0 or 1 items, instead of a list of all matching items
	if len(result.Items) == 0 {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) Save(user *domain.User) error {
	// First, the user struct is marshalled into a map that can be saved to the database
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Println("Error marshalling user", err)
		return err
	}
	// Then, the user is saved to the database using the PutItem method
	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      item,
	})
	if err != nil {
		log.Println("Error saving user: ", err)
	}
	return err
}
