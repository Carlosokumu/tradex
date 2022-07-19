package models

type PositionData struct {
	Data []struct {
		PositionId int    `json:"positionid,omitempty"`
		EntryPrice int    `json:"entryprice,omitempty"`
		TradeType  bool   `json:"tradetype,omitempty"`
		Quantity   string `json:quantity",omitempty"`
	}
}
