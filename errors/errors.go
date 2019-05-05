package errors

import "fmt"

type ProduceError struct {
	message string
}

func newProduceError(value string) error {
	return ProduceError{message: fmt.Sprintf("Request missing value: %s", value)}
}

func (c ProduceError) Error() string {
	return c.message
}
