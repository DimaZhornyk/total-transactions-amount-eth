package totalTransactionAmount

import (
	"strconv"
)

type RequestParams struct {
	BlockNumber int
}

type ExpectedParamsMissing struct{}

func (e *ExpectedParamsMissing) Error() string {
	return "Expected params are missing"
}

func NewRequestParams(vars map[string]string) (RequestParams, error) {
	var requestParams RequestParams
	if num, contains := vars["block_number"]; contains {
		blockNumber, err := strconv.Atoi(num)
		requestParams = RequestParams{BlockNumber: blockNumber}
		return requestParams, err
	}
	return requestParams, &ExpectedParamsMissing{}
}
