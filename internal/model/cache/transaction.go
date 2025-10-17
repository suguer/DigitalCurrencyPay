package cache

import (
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const (
	TransactionCacheKey = "transaction"
)

func TransactionCacheSet(transaction *mdb.Transaction) error {
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	toAddress := strings.ToUpper(transaction.ToAddress)
	// contactAddress := strings.ToUpper(transaction.ContractAddress)
	dao.Cache.Set(fmt.Sprintf("%s:%s", TransactionCacheKey, toAddress), jsonData)

	return nil
}
func TransactionCacheDelete(transaction *mdb.Transaction) error {
	toAddress := strings.ToUpper(transaction.ToAddress)
	return dao.Cache.Delete(fmt.Sprintf("%s:%s", TransactionCacheKey, toAddress))
}

func TransactionCacheGet(chain, toAddress, contactAddress string, transactionAmount float64) (*mdb.Transaction, error) {
	var transaction *mdb.Transaction
	toAddress = strings.ToUpper(toAddress)
	contactAddress = strings.ToUpper(contactAddress)
	matchAddress := []string{
		toAddress,
		fmt.Sprintf("41%s", toAddress),
		fmt.Sprintf("0X%s", toAddress),
	}
	if strings.HasPrefix(toAddress, "41") {
		matchAddress = append(matchAddress, fmt.Sprintf("0X%s", toAddress[2:]))
	}

	for _, address := range matchAddress {
		if data, err := dao.Cache.Get(fmt.Sprintf("%s:%s", TransactionCacheKey, address)); err == nil {
			if err := json.Unmarshal(data, &transaction); err != nil {
				continue
			}
			if matchTransaction(chain, contactAddress, transactionAmount, transaction) {
				return transaction, nil
			}
		}
	}
	return nil, fmt.Errorf("cache not found")
}

func matchTransaction(chain, contactAddress string, transactionAmount float64, transaction *mdb.Transaction) bool {
	if transaction.Chain != chain {
		return false
	}
	if !util.MatchAddress(transaction.ContractAddress, contactAddress) {
		return false
	}
	Amount := math.Abs(transaction.Amount)
	return Amount == transactionAmount

}
