package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Processor interface {
	ExportPayments(companyId string, paymentsId []string) (string, error)
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	paymentStorage db.PaymentRepository
	userStorage    db.UserRepository
}

func New(paymentRepo db.PaymentRepository, userRepo db.UserRepository) Processor {
	return &processor{
		paymentStorage: paymentRepo,
		userStorage:    userRepo,
	}
}

func (p *processor) ExportPayments(companyId string, paymentsId []string) (string, error) {
	var payments []domain.Payment

	payments, err := p.paymentStorage.GetPaymentsByMultipleIDs(paymentsId, companyId)
	if err != nil {
		return "", fmt.Errorf("error getting payments: %v", err)
	}

	// Access the S3 bucket to get the template
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1")},
	)
	if err != nil {
		return "", fmt.Errorf("error creating session: %v", err)
	}

	downloader := s3manager.NewDownloader(sess)

	tempFilePath := filepath.Join(os.TempDir(), "template.xlsx")

	file, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Descargar el archivo del bucket de S3
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("ASSET_S3_BUCKET_NAME")),
			Key:    aws.String(os.Getenv("TEMPLATE_FILE_PATH")),
		})
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}

	// Open the Excel file
	f, err := excelize.OpenFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	sheetName := "Template"

	// Iterate over the payments and add them to the Excel file
	for i, payment := range payments {
		rowIndex := i + 2

		statusStr := domain.ConfirmationStatusToString(payment.Status)

		if err != nil {
			return "", fmt.Errorf("error parsing confirmation status: %v", err)
		}
		formattedDate := payment.CreatedDate.Format("2006-01-02")

		_ = f.SetCellValue(sheetName, "A"+strconv.Itoa(rowIndex), payment.ID)
		_ = f.SetCellValue(sheetName, "B"+strconv.Itoa(rowIndex), payment.Amount)
		_ = f.SetCellValue(sheetName, "C"+strconv.Itoa(rowIndex), payment.ShopName)
		_ = f.SetCellValue(sheetName, "D"+strconv.Itoa(rowIndex), payment.Cuit)
		_ = f.SetCellValue(sheetName, "E"+strconv.Itoa(rowIndex), payment.Date)
		_ = f.SetCellValue(sheetName, "F"+strconv.Itoa(rowIndex), payment.Time)
		_ = f.SetCellValue(sheetName, "G"+strconv.Itoa(rowIndex), payment.Category)
		_ = f.SetCellValue(sheetName, "H"+strconv.Itoa(rowIndex), payment.ReceiptNumber)
		_ = f.SetCellValue(sheetName, "I"+strconv.Itoa(rowIndex), payment.ReceiptType)
		_ = f.SetCellValue(sheetName, "J"+strconv.Itoa(rowIndex), statusStr)
		_ = f.SetCellValue(sheetName, "K"+strconv.Itoa(rowIndex), formattedDate)
		_ = f.SetCellValue(sheetName, "L"+strconv.Itoa(rowIndex), payment.User.FullName)
		_ = f.SetCellValue(sheetName, "M"+strconv.Itoa(rowIndex), payment.User.Email)
	}

	// Save the file
	err = f.Save()
	if err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}

	// Upload the file to S3
	date := time.Now().Format("2006-01-02")
	fileName := date + ".xlsx"

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("ASSET_S3_BUCKET_NAME")),
		Key:    aws.String("exported-payments/" + fileName),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	// Return the presigned URL of the file using the aws generate presigned url function
	req, _ := s3.New(sess).GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("ASSET_S3_BUCKET_NAME")),
		Key:    aws.String("exported-payments/" + fileName),
	})
	url, err := req.Presign(3600 * time.Second)
	if err != nil {
		return "", fmt.Errorf("error presigning URL: %v", err)
	}

	// Update the payments Status to exported

	err = p.paymentStorage.UpdatePaymentsStatus(paymentsId, domain.Exported, companyId)
	if err != nil {
		return "", fmt.Errorf("error updating payments status: %v", err)
	}

	return url, nil
}

func (p *processor) ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error) {
	// Validate user
	isAuthorized, err := p.userStorage.IsUserAuthorized(email, companyId, allowedRoles)
	if err != nil {
		return false, err
	}
	if isAuthorized {
		return true, nil
	}
	return false, nil
}
