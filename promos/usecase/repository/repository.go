package repository

import (
	"context"

	"github.com/takadev15/sayakaya-promo-engine/promos/domain"
)

type IPromoRepository interface {
    InsertPromo(ctx context.Context, promo domain.Promo) error
    GetPromoTypes(ctx context.Context) ([]domain.PromoType, error)
}
