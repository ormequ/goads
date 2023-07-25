package app

import (
	"errors"
	"fmt"
)

var (
	ErrPermissionDenied = errors.New("this user cannot edit the ad")
	ErrNotFound         = errors.New("not found")
	ErrInvalidContent   = errors.New("invalid ad's content")
	ErrInvalidFilter    = errors.New("invalid ad's filter")
)

// Error is wrapped error (Err) with information: ID to check in what object of Type an error has happened
// and Details to add some additional info. Use errors.Is to check type of the error,
// e.g. errors.Is(Err, ErrPermissionDenied)
type Error struct {
	Err     error
	Type    string
	ID      int64
	Details string
}

func (e Error) Error() string {
	errText := fmt.Sprintf("error in %s", e.Type)
	if e.ID != -1 {
		errText = fmt.Sprintf("%s with ID %d", errText, e.ID)
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
