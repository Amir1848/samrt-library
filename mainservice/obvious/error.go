package obvious

import "context"

type BusinessError struct {
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

func ReturnPrefailedCondition(ctx context.Context, message string) error {
	return BusinessError{
		Message: message,
	}
}
