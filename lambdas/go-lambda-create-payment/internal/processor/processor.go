package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
	"go-lambda-create-payment/pkg/dto"
)

type Processor interface {
	CreatePayment(ctx context.Context, input *dto.CreatePaymentInput, email, companyId string) error
}

type processor struct {
	paymentStorage db.PaymentRepository
	userStorage    db.UserRepository
	teamStorage    db.TeamRepository
}

func New(paymentRepo db.PaymentRepository, userRepo db.UserRepository, teamRepo db.TeamRepository) Processor {
	return &processor{
		paymentStorage: paymentRepo,
		userStorage:    userRepo,
		teamStorage:    teamRepo,
	}
}

func (p *processor) CreatePayment(ctx context.Context, input *dto.CreatePaymentInput, email, companyId string) error {
	// Validate the status of the user (active or not)
	user, err := p.userStorage.GetUserIdByEmail(email, companyId)
	if err != nil {
		fmt.Println("Error getting user: ", err)
		return err
	}
	// Validate the purchase limit
	//purchaseLimit := user.PurchaseLimit
	//if float32(purchaseLimit) < input.Amount {
	//	return fmt.Errorf("the amount is greater than the purchase limit")
	//}
	//// Validate the monthly limit
	//monthlyLimit := user.MonthlyLimit
	//remainingMonthlyLimit := float32(monthlyLimit) - user.MonthlySpending
	//if remainingMonthlyLimit < input.Amount {
	//	return fmt.Errorf("Error: The amount is greater than the monthly limit")
	//}

	// Create payment
	payment, err := domain.NewPayment(input.Amount, input.ShopName, input.Cuit, input.Time, input.Category, input.ReceiptNumber, input.ReceiptType, input.ReceiptImageKey, user.ID, input.Date)
	if err != nil {
		fmt.Println("Error creating payment: ", err)
		return err
	}

	// Update the monthly spending of the user
	user.MonthlySpending += input.Amount
	if err := p.userStorage.UpdateUser(user, companyId); err != nil {
		fmt.Println("Error updating user: ", err)
		return err
	}
	// Save payment to db
	if err := p.paymentStorage.Save(payment, companyId); err != nil {
		fmt.Println("Error saving payment: ", err)
		return err
	}
	// If the user is part of a team, update the team's monthly spending
	userTeams, err := p.teamStorage.GetTeamByUserID(user.ID, companyId)
	fmt.Println("User teams: ", userTeams)
	if err != nil {
		fmt.Println("Error getting user teams ", err)
		return err
	}
	if len(userTeams) > 0 {
		fmt.Println("Updating team monthly spending")
		for _, userTeam := range userTeams {
			team, err := p.teamStorage.GetTeamByID(userTeam.TeamID, companyId)
			if err != nil {
				fmt.Println("Error getting team: ", err)
				return err
			}
			team.MonthlySpending += int(input.Amount)
			if err := p.teamStorage.UpdateTeamMonthlySpending(team.MonthlySpending, companyId); err != nil {
				fmt.Println("Error updating team monthly spending: ", err)
				return err
			}
		}
	}

	return nil
}
