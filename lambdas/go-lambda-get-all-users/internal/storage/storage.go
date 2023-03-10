package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func (r *UserRepository) GetAllUsers(filter string, attributeNames []string) ([]map[string]*dynamodb.AttributeValue, error) {
	// Set up the initial ScanInput with the table name and any other necessary parameters
	input := &dynamodb.ScanInput{
		TableName: aws.String("users"),
	}

	// Add filter expression to input if a filter was specified
	if filter != "" {
		filterExpression := ""
		expressionAttributeValues := map[string]*dynamodb.AttributeValue{}
		expressionAttributeNames := map[string]*string{}

		// Build the filter expression dynamically using the attribute names and filter value
		for i, attr := range attributeNames {
			if i > 0 {
				filterExpression += " OR "
			}
			filterExpression += "contains(#" + attr + ", :filter)"
			expressionAttributeValues[":filter"] = &dynamodb.AttributeValue{S: aws.String(filter)}
			expressionAttributeNames["#"+attr] = aws.String(attr)
		}

		input.FilterExpression = aws.String(filterExpression)
		input.ExpressionAttributeValues = expressionAttributeValues
		input.ExpressionAttributeNames = expressionAttributeNames
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
