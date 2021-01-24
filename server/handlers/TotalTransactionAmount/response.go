package totalTransactionAmount

type Response struct {
	Transactions int     `json:"transactions"`
	Amount       float64 `json:"amount"`
}
