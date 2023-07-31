package errwrap

import (
	"errors"
	"fmt"
)

// Error is wrapped error (Err) with information: ID to check in what object of Type an error has happened
// during Operation in Service. You can provide Details to add more information info
type Error struct {
	Err       error
	Operation string
	Service   string
	ID        int64
	Type      string
	Details   string
}

func (e Error) Error() string {
	errText := fmt.Sprintf("error during operation %s in service %s", e.Operation, e.Service)
	if e.ID != -1 {
		errText = fmt.Sprintf("%s in object with type %s and ID %d", errText, e.Type, e.ID)
	}
	errText = fmt.Sprintf("%s: %s", errText, e.Err.Error())
	if e.Details != "" {
		errText = fmt.Sprintf("%s (%s)", errText, e.Details)
	}
	return errText
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) WithDetails(details string) Error {
	e.Details = details
	return e
}

func (e Error) OnObject(objType string, id int64) Error {
	e.ID = id
	e.Type = objType
	return e
}

func New(err error, svc string, operation string) Error {
	return Error{
		Err:       err,
		ID:        -1,
		Service:   svc,
		Operation: operation,
		Type:      "unset",
	}
}

func JoinWithCaller(err error, caller string) error {
	if err == nil {
		return nil
	}
	return errors.Join(err, fmt.Errorf("called by %s", caller))
}
