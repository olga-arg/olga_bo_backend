package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func (r *UserRepository) GetAllUsers(filter string) ([]map[string]*dynamodb.AttributeValue, error) {
	// Set up the initial ScanInput with the table name and any other necessary parameters
	input := &dynamodb.ScanInput{
		TableName: aws.String("users"),
	}

	// Add filter expression to input if a filter was specified
	if filter != "" {
		input.FilterExpression = aws.String("#name = :filter")
		input.ExpressionAttributeNames = map[string]*string{
			"#name": aws.String("name"),
		}
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":filter": {
				S: aws.String(filter),
			},
		}
	}

	// Create a slice to hold all items retrieved from the database
	var allItems []map[string]*dynamodb.AttributeValue

	// Use a loop to retrieve all items, page by page
	for {
		// Call the Scan method with the current input to retrieve a page of items
		result, err := r.db.Scan(input)
		if err != nil {
			log.Println("Error getting users: ", err)
			return nil, err
		}

		// Append the items from the current page to the allItems slice
		allItems = append(allItems, result.Items...)

		// Check if there are more items to retrieve
		if result.LastEvaluatedKey == nil {
			// No more items to retrieve, break out of the loop
			break
		}

		// Set the ExclusiveStartKey parameter to the LastEvaluatedKey from the previous call to continue pagination
		input.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return allItems, nil
}
