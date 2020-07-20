package errors

type errorContext struct {
	Field   string
	Message string
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{
			errType:       customErr.errType,
			originalError: customErr.originalError,
			context:       context,
		}
	}

	return customError{
		errType:       Unknown,
		originalError: err,
		context:       context,
	}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {

		return map[string]string{
			"field":   customErr.context.Field,
			"message": customErr.context.Message,
		}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errType
	}

	return Unknown
}
