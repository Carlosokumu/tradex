package models

import "gorm.io/gorm"

//Struct will hold information about the various positions
type OpenPosition struct {
	gorm.Model
	PositionId int     `json:"PositionId,omitempty"`
	EntryPrice float32 `json:"EntryPrice,omitempty"`
	TradeType  string  `json:"TradeType,omitempty"`
	Quantity   float32 `json:"Quantity,omitempty"`
	EntryTime  string  `json:"EntryTime,omitempty"`
	TakeProfit float32 `json:"TakeProfit,omitempty"`
	StopLoss   float32 `json:"StopLoss,omitempty"`
}
