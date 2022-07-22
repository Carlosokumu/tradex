package models

import "gorm.io/gorm"

type PositionData struct {
	gorm.Model
	PositionId int     `json:"PositionId,omitempty"`
	EntryPrice float32 `json:"EntryPrice,omitempty"`
	TradeType  string  `json:"TradeType,omitempty"`
	Quantity   float32 `json:"Quantity,omitempty"`
	EntryTime  string  `json:"EntryTime,omitempty"`
}
