package storage

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go-lambda-create-user/pkg/domain"
)

type Storage interface {
	CreateUser(ctx context.Context, user *domain.User) error
}

type dynamoDB struct {
	client *dynamodb.DynamoDB
}

func NewDynamoDB() (Storage, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	client := dynamodb.New(sess)
	return &dynamoDB{
		client: client,
	}, nil
}

func (d *dynamoDB) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	item := map[string]*dynamodb.AttributeValue{
		"id": {
			S: &user.ID,
		},
		"name": {
			S: &user.Name,
		},
		"email": {
			S: &user.Email,
		},
	}

	_, err := d.client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("usersTable"),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}
