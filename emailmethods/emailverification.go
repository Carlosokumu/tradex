package emailmethods

import (
    "net/mail"
)

func Emailformatverifier (email string) (bool, error){

        _, err := mail.ParseAddress(email)
        if err != nil {
            return false,err
        } else {
            return true, nil
        }

}