package usecase

import (
	"context"

	"github.com/takadev15/sayakaya-promo-engine/promos/domain"
)

type IPromoUseCase interface {
	GenerateUserBirthdayPromo(ctx context.Context, user domain.User) error
}
