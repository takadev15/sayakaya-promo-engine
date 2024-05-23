package subscriber

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/takadev15/sayakaya-promo-engine/promos/domain"
	"github.com/takadev15/sayakaya-promo-engine/promos/usecase"
)

type redisSubscriberAdapter struct {
    promoUsecase usecase.IPromoUseCase
}

type RedisSubscriberAdapterOpts struct {
    PromoUsecase usecase.IPromoUseCase
}

func NewRedisAdapter(opts RedisSubscriberAdapterOpts) *redisSubscriberAdapter {
    return &redisSubscriberAdapter{
        promoUsecase: opts.PromoUsecase,
    }
}

func (adp redisSubscriberAdapter) SubscribeToBirthDay(client *redis.Client) {
    ctx := context.Background()
    subsciber := client.Subscribe(ctx, "birthdays_channel")
    user := User{}

    for {
        msg, err := subsciber.ReceiveMessage(ctx)
        if err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
            panic(err)
        }

        err = adp.promoUsecase.GenerateUserBirthdayPromo(ctx, domain.User{
            ID: user.ID,
            Name: user.Name,
            Email: user.Email,
            PhoneNumber: user.PhoneNumber,
        })
        if err != nil {
            panic(err)
        }
    }
}
