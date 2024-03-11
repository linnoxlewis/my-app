package models

import "errors"

var (
	ErrNotAvailableRooms   = errors.New("hotel room is not available for selected dates")
	ErrNotFoundInformation = errors.New("no availability information for selected date")
	ErrInternalServerError = errors.New("internal server error,please try later")
)
