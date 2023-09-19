package log

import "go.uber.org/zap"

var Logger zap.SugaredLogger

func Init() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Logger = *logger.Sugar()
	return nil
}
