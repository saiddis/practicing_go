package pkg

import "go.uber.org/zap"

func NewSugarLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}
