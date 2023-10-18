package controllers

import (
	"errors"
	"strings"
)

func namevalidator(name string) error {
	sname := strings.TrimSpace(name)
	if len(sname) < 3 || len(name) > 20 {
		return errors.New("name length must be between 4-20")
	}
	return nil
}
func numbervalidator(number string) error {
	snumber := strings.TrimSpace(number)
	if len(snumber) != 10 {
		return errors.New("mobile number must be 10 digits")
	}
	return nil
}
func emailvalidator(email string) error {
	semail := strings.TrimSpace(email)
	if !strings.Contains(semail, "@") {
		return errors.New("invalid Email Id")
	}
	if len(semail) < 4 {
		return errors.New("invalid Email Id")
	}
	return nil
}
func passwordvalidator(password string) error {
	spassword := strings.TrimSpace(password)
	if len(spassword) < 5 {
		return errors.New("password at least 5 character")
	}
	return nil
}
