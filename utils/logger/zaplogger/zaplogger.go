package zaplogger

import "go.uber.org/zap"

func New(production bool) *zap.SugaredLogger {
	var conf zap.Config
	if production {
		conf = zap.NewProductionConfig()
	} else {
		conf = zap.NewDevelopmentConfig()
	}

	conf.EncoderConfig.EncodeLevel = capitalLevelEncoder
	conf.DisableStacktrace = true
	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}
	return logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}
