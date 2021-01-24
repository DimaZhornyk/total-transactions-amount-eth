package tests

import (
	"math/big"
	"testing"
	"transaction-amount-eth/server/models"
	"transaction-amount-eth/server/utils"
)

func TestCreateRequestWithParams(t *testing.T) {
	manager := models.EtherScanManager{Url: "https://api.etherscan.io/api", Module: "proxy", ApiKey: "test_api_key"}
	tables := []struct {
		action      string
		blockNumber int
		result      string
	}{
		{"eth_blockNumber", 123, "https://api.etherscan.io/api?action=eth_blockNumber&apikey=test_api_key&boolean=true&module=proxy&tag=7b"},
		{"getblockreward", 456, "https://api.etherscan.io/api?action=getblockreward&apikey=test_api_key&boolean=true&module=proxy&tag=1c8"},
		{"eth_getBlockByNumber", 6789876, "https://api.etherscan.io/api?action=eth_getBlockByNumber&apikey=test_api_key&boolean=true&module=proxy&tag=679af4"},
	}
	for _, table := range tables {
		params := map[string]string{
			"module":  manager.Module,
			"action":  table.action,
			"boolean": "true",
			"tag":     utils.Dec2hex(table.blockNumber),
			"apikey":  manager.ApiKey,
		}
		uri := utils.CreateRequestUrlWithParams(manager.Url, params)
		if uri != table.result {
			t.Errorf("Request with params is incorrect with input data: action(%s), blockNumber(%d), got: %s, want: %s.", table.action, table.blockNumber, uri, table.result)
		}
	}
}

func TestDec2Hex(t *testing.T) {
	tables := []struct {
		dec int
		hex string
	}{
		{100, "64"},
		{123123123, "756b5b3"},
		{7777, "1e61"},
	}
	for _, table := range tables {
		res := utils.Dec2hex(table.dec)
		if res != table.hex {
			t.Errorf("Dec2Hex returned wrong result with input data: dec(%d), got: %s, want: %s.", table.dec, res, table.hex)
		}
	}
}

func TestHex2BigInt(t *testing.T) {
	tables := []struct {
		hex string
		bigInt *big.Int
	}{
		{"0x64", big.NewInt(100)},
		{"0x756b5b3", big.NewInt(123123123)},
		{"0x1e61", big.NewInt(7777)},
	}
	for _, table := range tables {
		res := utils.Hex2BigInt(table.hex)
		if res.Cmp(table.bigInt) != 0 {
			t.Errorf("Hex2BigInt returned wrong result with input data: hex(%s), got: %d, want: %d.", table.hex, res, table.bigInt)
		}
	}
}
