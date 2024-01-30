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
	"go-lambda-me/internal/processor"
	"net/http"
	"os"
	"strings"
)

type MeHandler struct {
	processor processor.Processor
}

func NewMeHandler(p processor.Processor) *MeHandler {
	return &MeHandler{processor: p}
}

func (h *MeHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authHeader := request.Headers["authorization"]

	splitAuthHeaders := strings.Split(authHeader, " ")
	if len(splitAuthHeaders) != 2 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing or invalid authorization header, actual Header: " + authHeader,
		}, nil
	}
	userPoolId := os.Getenv("USER_POOL_ID")
	pubKeyUrl := fmt.Sprintf("https://cognito-idp.us-east-1.amazonaws.com/%s/.well-known/jwks.json", userPoolId)
	keySet, err := jwk.Fetch(context.Background(), pubKeyUrl)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	token, err := jwt.Parse([]byte(splitAuthHeaders[1]), jwt.WithKeySet(keySet), jwt.WithValidate(true))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}
	username, _ := token.Get("username")
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))

	// Create a new Cognito Identity Provider client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Prepare the input for the AdminGetUser call
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(username.(string)),
	}

	// Call AdminGetUser to retrieve user data
	result, err := cognitoClient.AdminGetUser(input)
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

	userInformation, err := h.processor.GetUserInformation(context.Background(), email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(userInformation)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
