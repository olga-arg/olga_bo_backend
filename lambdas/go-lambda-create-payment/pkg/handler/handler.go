package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go-lambda-create-payment/internal/processor"
	"go-lambda-create-payment/pkg/dto"
	"net/http"
	"os"
	"strings"
)

type CreatePaymentHandler struct {
	processor processor.Processor
}

func NewCreatePaymentHandler(p processor.Processor) *CreatePaymentHandler {
	return &CreatePaymentHandler{processor: p}
}

func (h *CreatePaymentHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader := request.Headers["authorization"]
	print("authHeader: ", authHeader)
	splitAuthHeaders := strings.Split(authHeader, " ")
	if len(splitAuthHeaders) != 2 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing or invalid authorization header, actual Header: " + authHeader,
		}, nil
	}
	userPoolId := os.Getenv("USER_POOL_ID")
	print("upi: ", userPoolId)
	pubKeyUrl := fmt.Sprintf("https://cognito-idp.us-east-1.amazonaws.com/%s/.well-known/jwks.json", userPoolId)
	print("pb: ", pubKeyUrl)
	keySet, err := jwk.Fetch(context.Background(), pubKeyUrl)
	print("kS: ", keySet)
	if err != nil {
		print("0 err: ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	token, err := jwt.Parse([]byte(splitAuthHeaders[1]), jwt.WithKeySet(keySet), jwt.WithValidate(true))
	print("token: ", token)
	if err != nil {
		print("1 err: ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}
	username, _ := token.Get("username")
	print("username: ", username)
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))

	// Create a new Cognito Identity Provider client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Prepare the input for the AdminGetUser call
	adminGetUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(username.(string)),
	}

	// Call AdminGetUser to retrieve user data
	result, err := cognitoClient.AdminGetUser(adminGetUserInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	// access the user attributes using result.UserAttributes
	// Get the email attribute:
	email := ""
	for _, attr := range result.UserAttributes {
		if *attr.Name == "email" {
			email = *attr.Value
			break
		}
	}
	print("email: ", email)
	var input dto.CreatePaymentInput

	// Unmarshal the request body into the input struct
	err = json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		print(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreatePayment(context.Background(), &input, email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "Payment created successfully",
	}, nil
}
