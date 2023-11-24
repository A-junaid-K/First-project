package controllers

import (
	"errors"
	"regexp"
	"strings"
)

func namevalidator(name string) error {
	var sname string
	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if nameRegex.MatchString(name) {
		sname = name
	} else {
		return errors.New("name is not valid")
	}
	sname = strings.TrimSpace(name)
	if len(sname) < 2 || len(name) > 20 {
		return errors.New("name length must be between 3-20")
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
	if !strings.Contains(semail, "@") || !strings.HasSuffix(semail, ".com") {
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
