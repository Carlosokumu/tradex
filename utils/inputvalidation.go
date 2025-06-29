package utils

import (
	"errors"

	"github.com/carlosokumu/dubbedapi/dtos"
	"github.com/carlosokumu/dubbedapi/emailmethods"
	"github.com/carlosokumu/dubbedapi/stringmethods"
)

func ValidateUserInput(user *dtos.UserDto) error {
	if usernamelength := len(user.UserName); usernamelength < 4 {
		return errors.New("username should be more of 4 or more characters")
	}
	if emailformatcredibility, _ := emailmethods.Emailformatverifier(user.Email); !emailformatcredibility {
		return errors.New("please check that your email is correctly formatted")
	}
	if passwordlength := stringmethods.Charactercount(user.Password); passwordlength < 8 {
		return errors.New("a password should be of 8 or more characters")
	}
	return nil
}
