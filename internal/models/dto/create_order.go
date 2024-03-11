package dto

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	emailReg           = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	errCannotBeEmpty   = "%s cannot be empty"
	errIncorrectFormat = "%s has incorrect format"
)

type Order struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (o *Order) Validate() error {
	var errs []string

	if o.HotelID == "" {
		errs = append(errs, fmt.Sprintf(errCannotBeEmpty, "hotel id"))
	}

	if o.RoomID == "" {
		errs = append(errs, fmt.Sprintf(errCannotBeEmpty, "room id"))
	}

	if o.UserEmail == "" {
		errs = append(errs, fmt.Sprintf(errCannotBeEmpty, "user email"))
	} else if !emailReg.MatchString(o.UserEmail) {
		errs = append(errs, fmt.Sprintf(errIncorrectFormat, "user email"))
	}

	if o.From.IsZero() {
		errs = append(errs, fmt.Sprintf(errCannotBeEmpty, "from"))
	}

	if o.To.IsZero() {
		errs = append(errs, fmt.Sprintf(errCannotBeEmpty, "to"))
	}

	if o.From.After(o.To) {
		errs = append(errs, fmt.Sprintf(errIncorrectFormat, "from-to"))
	}

	if len(errs) > 0 {
		strErr := strings.Join(errs, "\n ")
		return errors.New(strErr)
	}

	return nil
}
