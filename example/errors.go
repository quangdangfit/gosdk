package main

import (
	"github.com/quangdangfit/gocommon/errors"
	"github.com/quangdangfit/gocommon/utils/logger"
)

func main() {
	err := errors.New("this is error", true)

	errWithContext := errors.AddErrorContext(err, "field", "message")

	if err != nil {
		logger.Error(err.Error())
	}

	if errWithContext != nil {
		logger.Error(errors.GetType(errWithContext))
	}

	err = errors.New("an_error", true)
	wrappedError := errors.BadRequest.Wrapf(err, "bad request %s", "not found")

	logger.Info(errors.GetType(wrappedError))
	logger.Info(wrappedError.Error())
}
