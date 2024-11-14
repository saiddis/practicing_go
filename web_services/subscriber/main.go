package main

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		sugarLogger.Fatalf("error reading config file: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_host") + ":" + viper.GetString("redis_port"),
		Password: "",
		DB:       0,
	})

	nc, err := nats.Connect(viper.GetString("nats_addr"))
	if err != nil {
		sugarLogger.Fatalf("error connecting to NATS server: %v", err)
	}
	defer nc.Drain()

	var name string
	var greet string
	ctx := context.Background()
	_, err = nc.Subscribe("chat.*", func(msg *nats.Msg) {
		name = msg.Subject[5:]
		greet = string(msg.Data)
		msg.Respond([]byte(msg.Subject))

		err = rdb.Set(ctx, name, greet, 1*time.Minute).Err()
		if err != nil {
			sugarLogger.Errorf("error setting values to redis: %v", err)
			return
		}
		val, err := rdb.Get(ctx, name).Result()
		if err != nil {
			sugarLogger.Errorf("error getting values from redis: %v", err)
			return
		}
		sugarLogger.Infof("%s: %s", name, val)
	})

	if err != nil {
		sugarLogger.Fatalf("error subscribing to subject: %v", err)
	}

	select {}
}
