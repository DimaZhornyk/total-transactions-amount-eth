package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"transaction-amount-eth/server/db"
	"transaction-amount-eth/server/handlers"
	totalTransactionAmount "transaction-amount-eth/server/handlers/TotalTransactionAmount"
	"transaction-amount-eth/server/models"
)

type MockDbManager struct{}

func (m *MockDbManager) GetBlockInfo(blockNumber int) (error, bool, db.BlockInfo) {
	return nil, false, db.BlockInfo{}
}

func (m *MockDbManager) InsertBlockInfo(blockNumber int , transactionsNumber int, totalEth float64) error {
	return nil
}

func TestHttp(t *testing.T) {
	tables := []struct {
		blockNumber  int
		transactions int
		total        float64
	}{
		{11509797, 155, 2.285405},
		{11508993, 241, 1130.987085},
		{109789, 1, 4.99877},
	}
	h := handlers.BaseHandler{
		DbManager:        &MockDbManager{},
		EtherScanManager: models.EtherScanManager{Url: "https://api.etherscan.io/api", Module: "proxy", ApiKey: "HWQNPRTV7GFNI2WEYMWVKIWUUYAHWXHY8V"},
	}
	router := mux.NewRouter()
	ttaHandler := totalTransactionAmount.NewHandler(&h)
	router.HandleFunc("/api/block/{block_number:[0-9]+}/total", ttaHandler.TotalTransactionAmount).Methods("GET")

	srv := httptest.NewServer(router)
	defer srv.Close()
	for _, table := range tables {
		res, err := http.Get(fmt.Sprintf("%s/api/block/%d/total", srv.URL, table.blockNumber))
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("status not OK")
		}

		body, err := ioutil.ReadAll(res.Body)
		var result totalTransactionAmount.Response
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result.Transactions != table.transactions || math.Abs(result.Amount-table.total) > 0.00001 {
			t.Errorf("HttpTest failed, returned wrong result with input data: blockNumber(%d), got: transactions(%d),"+
				" totalAmount(%f), want: transactions(%d), totalAmount(%f).",
				table.blockNumber, result.Transactions, result.Amount, table.transactions, table.total)
		}
		res.Body.Close()
	}
}
