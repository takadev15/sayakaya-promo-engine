package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/takadev15/sayakaya-promo-engine/promos/domain"
	sqlDB "github.com/takadev15/sayakaya-promo-engine/promos/external/database/sql"
)

type promoRepository struct {
    sqlDB sqlDB.SQLDatabase
}

type PromoRepositoryOpts struct {
    SQLDB sqlDB.SQLDatabase
}

func NewPromoRepository(opts PromoRepositoryOpts) *promoRepository {
    return &promoRepository{
        sqlDB: opts.SQLDB,
    }
}

func (repo promoRepository) InsertPromo(ctx context.Context, promo domain.Promo) error {
    _, err := repo.sqlDB.Exec(
        ctx, queryInsertBirthdayPromo,
        promo.PromoCode,
        promo.UserID,
        promo.PromoTypeID,
        promo.ValidUntil,
        promo.Status,
    )
    if err != nil {
        return err
    }
    return nil
}

func (repo promoRepository) GetPromoTypes(ctx context.Context) ([]domain.PromoType, error) {
    var queryResult []PromoType
    var results []domain.PromoType 

    err := repo.sqlDB.Get(ctx, &queryResult, querySelectPromoTypes)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("Promo type empty")
            return nil, nil
        }
        return nil, err
    }

    for _, res := range(queryResult) {
        results = append(results, domain.PromoType{
            ID: res.PromoTypeID,
            Name: res.Name,
            Rule: res.Rule,
        })
    }

    return results, nil
}
