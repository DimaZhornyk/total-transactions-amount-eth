package handlers

import (
	"transaction-amount-eth/server/db"
	"transaction-amount-eth/server/models"
)

type BaseDbManager interface {
	GetBlockInfo(int) (error, bool, db.BlockInfo)
	InsertBlockInfo(int, int, float64) error
}

type BaseHandler struct {
	DbManager        BaseDbManager
	EtherScanManager models.EtherScanManager
}
