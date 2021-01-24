package totalTransactionAmount

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"transaction-amount-eth/server/handlers"
	"transaction-amount-eth/server/models"
	"transaction-amount-eth/server/utils"
)

type Handler struct {
	*handlers.BaseHandler
}

func NewHandler(bh *handlers.BaseHandler) Handler {
	return Handler{bh}
}

func (handler *Handler) TotalTransactionAmount(w http.ResponseWriter, r *http.Request) {
	requestParams, err := NewRequestParams(mux.Vars(r))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err, exists, blockInfo := handler.DbManager.GetBlockInfo(requestParams.BlockNumber)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		response := Response{
			Transactions: blockInfo.Transactions,
			Amount:       blockInfo.Total,
		}
		utils.JsonResponse(w, response, http.StatusOK)
	} else {
		res, err := handler.EtherScanManager.FetchBlockData("eth_getBlockByNumber", requestParams.BlockNumber)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error accessing remote server", http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		var result models.EtherScanData
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		totalWei := result.CountTotalWeiQuantity()
		totalTransactions := result.CountTransactionsQuantity()
		totalEth, _ := totalWei.Quo(totalWei, big.NewFloat(math.Pow(10, 18))).Float64()
		err = handler.DbManager.InsertBlockInfo(requestParams.BlockNumber, totalTransactions, totalEth)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error getting data ", http.StatusInternalServerError)
			return
		}
		response := Response{
			Transactions: totalTransactions,
			Amount:       totalEth,
		}
		utils.JsonResponse(w, response, http.StatusOK)
	}
}
