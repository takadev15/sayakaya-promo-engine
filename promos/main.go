package main

import (
	"github.com/takadev15/sayakaya-promo-engine/promos/adapter/subscriber"
	"github.com/takadev15/sayakaya-promo-engine/promos/config"
	"github.com/takadev15/sayakaya-promo-engine/promos/external/database/sql"
	"github.com/takadev15/sayakaya-promo-engine/promos/external/pubsub"
	"github.com/takadev15/sayakaya-promo-engine/promos/usecase"
	"github.com/takadev15/sayakaya-promo-engine/promos/usecase/repository"
)

func main() {
    cfg := config.LoadConfig()

    db, err := sql.NewPostgresClient(sql.ConnectionOption{
        Host: cfg.PSQL.Host,
        Port: cfg.PSQL.Port,
        User: cfg.PSQL.User,
        Password: cfg.PSQL.Password,
        Database: cfg.PSQL.Database,
    })
    if err != nil {
        panic(err)
    }

    redisC := pubsub.NewRedisClient("6379")

    promoRepository := repository.NewPromoRepository(repository.PromoRepositoryOpts{
        SQLDB: db,
    })

    promoUsecase := usecase.NewPromoUsecase(usecase.PromoUsecaseOpts{
        PromoRepository: promoRepository,
    })

    pubSubAdapter := subscriber.NewRedisAdapter(subscriber.RedisSubscriberAdapterOpts{
        PromoUsecase: promoUsecase,
    })

    pubSubAdapter.SubscribeToBirthDay(redisC)

}
