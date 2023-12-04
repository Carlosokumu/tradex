package models

type  AccessRefreshaccountsecret struct {
	AccessToken string    `form:"access_token"`
	RefreshToken string `form:"refresh_token"`
	Client_id string    `form:"client_id"`
	Secret string `form:"secret"`
}