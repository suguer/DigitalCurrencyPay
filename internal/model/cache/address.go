package cache

import (
	"DigitalCurrency/internal/model/dao"
	"context"
	"fmt"
	"strings"
	"time"
)

func GetInUseAddress(chain string, amount float64) []string {
	keys, _ := dao.Rdb.Keys(context.Background(), fmt.Sprintf("InUseAddress:%s:*", chain)).Result()
	fmt.Printf("keys: %v\n", keys)
	data := make([]string, 0)
	for _, key := range keys {
		arr := strings.Split(key, ":")
		if arr[1] != chain {
			continue
		}
		data = append(data, arr[2])
	}
	return data
}

func AddInUseAddress(chain string, address string, amount float64) {
	dao.Rdb.Set(context.Background(), fmt.Sprintf("InUseAddress:%s:%s", chain, address), amount, time.Minute*30)
}
