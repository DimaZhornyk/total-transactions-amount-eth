package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func JsonResponse(w http.ResponseWriter, data interface{}, c int) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	_, err = fmt.Fprintf(w, "%s", j)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func Dec2hex(dec int) string {
	return strconv.FormatInt(int64(dec), 16)
}

func Hex2BigInt(hex string) *big.Int {
	i := new(big.Int)
	if strings.HasPrefix(hex, "0x") {
		i.SetString(hex[2:], 16)
	} else {
		i.SetString(hex, 16)
	}
	return i
}

func CreateRequestUrlWithParams(url string, args map[string]string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	q := req.URL.Query()
	for k, v := range args {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}
