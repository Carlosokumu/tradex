package dtos

type TradingAccountDto struct {
	Platform  string `json:"platform"`
	AccountId string `json:"account_id"`
	UserName  string `json:"username"`
}
