package main

import (
	"time"

	"github.com/nats-io/nats.go"
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

	nc, _ := nats.Connect(viper.GetString("nats_addr"))
	defer nc.Drain()
	_, err = nc.Subscribe("greet.*", func(msg *nats.Msg) {
		//name := msg.Subject[6:]
		msg.Respond([]byte(msg.Subject))
	})
	if err != nil {
		sugarLogger.Fatalf("error subscribing to subject: %v", err)
	}

	select {}
	//PublishSubscribe(nc, sugarLogger)

	//rdb := redis.NewClient(&redis.Options{
	//	Addr: viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT"),
	//})

}

func PublishSubscribe(nc *nats.Conn, logger *zap.SugaredLogger) {
	defer nc.Drain()
	sub, _ := nc.Subscribe("greet.*", func(msg *nats.Msg) {
		name := msg.Subject[6:]
		msg.Respond([]byte("hello, " + name))
	})

	rep, err := nc.Request("greet.joe", nil, time.Second)
	if err != nil {
		logger.Logf(1, "error making request: %v", err)
	}
	logger.Logf(1, "Recieved reply: %v", string(rep.Data))

	rep, err = nc.Request("greet.sue", nil, time.Second)
	if err != nil {
		logger.Logf(1, "error making request: %v", err)
	}
	logger.Logf(1, "Recieved reply: %s", string(rep.Data))

	rep, err = nc.Request("greet.bob", nil, time.Second)
	if err != nil {
		logger.Logf(1, "error making request: %v", err)
	}
	logger.Logf(1, "Recieved reply: %s", string(rep.Data))

	sub.Unsubscribe()

	_, err = nc.Request("greet.joe", nil, time.Second)
	if err != nil {
		logger.Logf(1, "error making request: %v", err)
	}
	logger.Logf(1, "error requesting service2: %v", err)

}
