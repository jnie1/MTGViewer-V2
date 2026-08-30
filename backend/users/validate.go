package users

import (
	"errors"
	"strings"
)

var EmptyUsername = errors.New("username is empty")
var PasswordShort = errors.New("password is too short")

func Validate(req SignupRequest) error {
	username := strings.TrimSpace(req.Username)
	if len(username) == 0 {
		return EmptyUsername
	}

	if len(req.Password) < 8 {
		return PasswordShort
	}

	return nil
}
