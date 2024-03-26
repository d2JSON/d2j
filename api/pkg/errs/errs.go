package errs

import "strings"

// Error implements the Error interface with error marshaling.
type Error struct {
	Message string `json:"message"`
}

// New creates new Error.
func New(message string) *Error {
	return &Error{Message: message}
}

func (e Error) Error() string {
	return e.Message
}

// IsExpected checks whether the input error can be
func IsExpected(err error) bool {
	_, ok := err.(*Error) //nolint:errorlint
	return ok
}

// HasAnyGivenMessage returns true if given error has any given message.
// You can pass infinite number of messages into this function, but at least one is always required
func HasAnyGivenMessage(err error, firstMsg string, otherMsgs ...string) bool {
	if strings.Contains(err.Error(), firstMsg) {
		return true
	}
	for _, msg := range otherMsgs {
		if strings.Contains(err.Error(), msg) {
			return true
		}
	}
	return false
}
