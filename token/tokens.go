package token

import "github.com/carlosokumu/dubbedapi/models"

type Token struct {
	User        models.User `json:"user"`
	TokenString string   `json:"token"`
}
