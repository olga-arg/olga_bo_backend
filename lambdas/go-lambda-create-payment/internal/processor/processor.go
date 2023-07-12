package processor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go-lambda-create-payment/internal/storage"
	"go-lambda-create-payment/pkg/domain"
	"go-lambda-create-payment/pkg/dto"
	"os"
)

type Processor interface {
	CreatePayment(ctx context.Context, input *dto.CreatePaymentInput) error
}

type processor struct {
	storage     storage.PaymentRepository
	teamStorage storage.TeamRepository
}

func New(paymentRepo storage.PaymentRepository, teamRepo storage.TeamRepository) Processor {
	return &processor{
		storage:     paymentRepo,
		teamStorage: teamRepo,
	}
}

func (p *processor) CreatePayment(ctx context.Context, input *dto.CreatePaymentInput) error {
	//shopName := "Amazon"
	//cardID := "1234567890"
	//userID := "1234567890"
	//category := "Groceries"
	//receipt := "https://s3.amazonaws.com/your-bucket-name/receipt.jpg"
	//paymentType := domain.Card

	//if input.Amount < 0 {
	//	fmt.Println("Error: amount cannot be less than 0")
	//	return errors.New("amount cannot be less than 0")
	//}

	//// If there is an input.Label then it is a payment to a team
	//if input.Label != "" {
	//	// Find team by ID
	//	team, err := p.teamStorage.FindTeamByID(input.Label)
	//	if err != nil {
	//		fmt.Println("Error finding team: ", err)
	//		return err
	//	}
	//	// If the team exists then set the payment type to team
	//	err = p.teamStorage.UpdateTeamMonthlySpending(team, input.Amount)
	//	if err != nil {
	//		fmt.Println("Error updating payment type to team: ", err)
	//		return err
	//	}
	//}

	// Creates a new payment - When receiving payments via Bank API, use the update payment function to upload the receipt from the app.
	payment, err := domain.NewPayment(input.Amount, input.ShopName, input.CardID, input.UserID, input.Category, input.Receipt)
	if err != nil {
		fmt.Println("Error creating payment: ", err)
		return err
	}
	if input.Receipt != "" {
		// base64 decode the receipt
		file, err := base64.StdEncoding.DecodeString(input.Receipt)
		if err != nil {
			fmt.Println("Error decoding receipt: ", err)
		}
		// Upload the receipt to S3
		// obtain the KeyID and SecretKey from the environment variables
		ac := os.Getenv("S3_USER_AC")
		sac := os.Getenv("S3_USER_SAC")
		s3Config := &aws.Config{
			Region:      aws.String("sa-east-1"),
			Credentials: credentials.NewStaticCredentials(ac, sac, ""),
		}
		s3Session, err := session.NewSession(s3Config)
		if err != nil {
			fmt.Println("Error creating session: ", err)
		}
		uploader := s3manager.NewUploader(s3Session)
		inputAws := &s3manager.UploadInput{
			Bucket:      aws.String("prod-olga-backend-receipts"), // bucket's name
			Key:         aws.String(payment.ID),                   // files destination location
			Body:        bytes.NewReader(file),                    // content of the file
			ContentType: aws.String("image/jpg"),                  // content type
			ACL:         aws.String("public-read"),
		}
		output, err := uploader.UploadWithContext(context.Background(), inputAws)
		if err != nil {
			fmt.Println("Error uploading file: ", err)
		} else {
			payment.Receipt = output.Location
		}
	}
	// Saves the payment to the database if it doesn't already exist
	if err := p.storage.Save(payment); err != nil {
		fmt.Println("Error saving payment: ", err)
		return err
	}
	// Returns
	return nil
}
