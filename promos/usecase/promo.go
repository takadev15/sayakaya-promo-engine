package usecase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/takadev15/sayakaya-promo-engine/promos/domain"
	"github.com/takadev15/sayakaya-promo-engine/promos/usecase/repository"
)

type promoUsecase struct {
	promoRepository repository.IPromoRepository
}

type PromoUsecaseOpts struct {
	PromoRepository repository.IPromoRepository
}

func NewPromoUsecase(opts PromoUsecaseOpts) *promoUsecase {
	return &promoUsecase{
		promoRepository: opts.PromoRepository,
	}
}

func (u promoUsecase) GenerateUserBirthdayPromo(ctx context.Context, user domain.User) error {
	promoTypes, err := u.promoRepository.GetPromoTypes(ctx)
	if err != nil {
		return err
	}

	if len(promoTypes) == 0 {
		log.Println("No promo types in databases")
		return fmt.Errorf("no promo types in database")
	}

	// Select random promo types
	randomIndex := rand.Intn(len(promoTypes))
	choosedType := promoTypes[randomIndex]

	// Generate Promo Code
	codePrefix := "SYKY-"
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(90000) + 10000
	promoCode := fmt.Sprintf("%s%d", codePrefix, randomNumber)

	// Set expiration date to today late night
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location())

	promo := domain.Promo{
		UserID:      user.ID,
		PromoTypeID: choosedType.ID,
		Status:      0,
		PromoCode:   promoCode,
		ValidUntil:  endOfDay,
	}

	err = u.promoRepository.InsertPromo(ctx, promo)
	if err != nil {
		log.Println("Error while querying data")
		return err
	}

	// send email and whatsapp concurrently message using goroutine
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- u.sendEmail(user.Email, "Send email")
	}()
	go func() {
		errChannel <- u.sendWhatsapp(user.PhoneNumber, "send message")
	}()
	// Collect errors from both goroutines
	for i := 0; i < 2; i++ {
		err := <-errChannel
		if err != nil {
			fmt.Println("Error occurred:", err)
			return err
		} else {
			fmt.Println("Task completed successfully")
		}
	}

	return nil
}

func (u promoUsecase) sendEmail(email string, message string) error {
	fmt.Printf("sending mail to %s with message %s", email, message)
	return nil
}

func (u promoUsecase) sendWhatsapp(phoneNumber string, message string) error {
	fmt.Printf("sending message to %s with message %s", phoneNumber, message)
	return nil
}
