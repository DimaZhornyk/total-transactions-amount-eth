package models

import (
	"math/big"
	"net/http"
	"transaction-amount-eth/server/utils"
)

type EtherScanData struct {
	Result struct {
		Transactions []struct{
			Value string `json:"value"`
		} `json:"transactions"`
	} `json:"result"`
}

type EtherScanManager struct {
	Url    string
	Module string
	ApiKey string
}

func (m *EtherScanManager) FetchBlockData(action string, blockNumber int) (*http.Response, error) {
	uri := utils.CreateRequestUrlWithParams(m.Url, map[string]string{
		"module":  m.Module,
		"action":  action,
		"boolean": "true",
		"tag":     utils.Dec2hex(blockNumber),
		"apikey":  m.ApiKey,
	})
	return http.Get(uri)
}

func (d *EtherScanData) CountTotalWeiQuantity() *big.Float {
	totalWei := big.NewInt(0)
	for _, transaction := range d.Result.Transactions {
		totalWei.Add(totalWei, utils.Hex2BigInt(transaction.Value))
	}
	return new(big.Float).SetInt(totalWei)
}

func (d *EtherScanData) CountTransactionsQuantity() int {
	return len(d.Result.Transactions)
}
