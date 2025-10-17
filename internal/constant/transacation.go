package constant

import "DigitalCurrency/internal/model/mdb"

func TransactionStatusMap(status int) string {
	m := map[int]string{
		mdb.TransactionStatusInit:       "unpaid",
		mdb.TransactionStatusSuccess:    "success",
		mdb.TransactionStatusFail:       "failed",
		mdb.TransactionStatusCollecting: "collecting",
		mdb.TransactionStatusCollected:  "collected",
	}
	return m[status]
}
