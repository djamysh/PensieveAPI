package utils

type EventError struct {
	Message string
}

func (err *EventError) Error() string {
	return err.Message
}
