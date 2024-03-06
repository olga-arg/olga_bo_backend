package services

// We need to create a function that will recieve a var validUsers []dto.CreateUserInput
// and for each user in the array, we need to create a new user in the database, cognito and send an email to the user
// But, we need to do this in a SQS queue, so we need to create a new service that will handle the queue

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go-lambda-bulk-create-users/pkg/dto"
	"os"
	"time"
)

func getSqsQueue(queueName string) (err error, queueURL string) {
	// Create a new SQS service
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	// Get the URL for the queue
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Println("Error getting queue URL for", queueName, ":", err)
		return err, ""
	}

	return err, *resultURL.QueueUrl
}

func SendMessageToQueue(validUsers []dto.CreateUserInput, companyId string) error {
	queueName := os.Getenv("SQS_QUEUE_NAME")
	// Obtener la URL de la cola
	err, queueURL := getSqsQueue(queueName)
	if err != nil {
		return err
	}

	fmt.Println("Queue URL:", queueURL)

	// Crear un nuevo servicio SQS
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	// Agregar los diccionarios de usuarios a la lista
	for _, user := range validUsers {
		time1 := time.Now()
		userDict := map[string]interface{}{
			"company_id": companyId,
			"name":       user.Name,
			"surname":    user.Surname,
			"email":      user.Email,
			"role":       user.Role,
		}

		// Convertir la lista de diccionarios a JSON
		messageBody, err := json.Marshal(userDict)
		if err != nil {
			fmt.Println("Error marshalling message body:", err)
			return err
		}

		fmt.Println("Sending message to queue:", string(messageBody))

		// Enviar el mensaje a la cola
		_, err = svc.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(string(messageBody)),
			QueueUrl:    &queueURL,
		})
		if err != nil {
			fmt.Println("Error sending message to queue:", err)
			return err
		}
		fmt.Println("Message sent successfully to queue")
		fmt.Println("Time: ", time.Since(time1))
	}

	return nil
}
