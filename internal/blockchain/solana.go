package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

type Solana struct {
	Blockchain
	rpcClient *rpc.Client
}

func NewSolana(ctx context.Context, conf *config.EthConfig) *Solana {
	rpcClient := NewRPC(conf.GrpcAddress)
	return &Solana{
		Blockchain: Blockchain{
			Chain:  conf.Name,
			Config: conf,
			ctx:    ctx,
		},
		rpcClient: rpcClient,
	}
}

func (s *Solana) GetBalance(addr string) (uint64, error) {
	solanaAddr, err := solana.PublicKeyFromBase58(addr)
	if err != nil {
		return 0, err
	}
	out, err := s.rpcClient.GetBalance(s.ctx, solanaAddr, rpc.CommitmentConfirmed)
	if err != nil {
		return 0, err
	}
	return out.Value, nil
}
func (s *Solana) GetNowBlockId(blockID ...int) (int, time.Time, error) {
	blockhash, err := s.rpcClient.GetLatestBlockhash(s.ctx, rpc.CommitmentFinalized)
	if err != nil {
		return 0, time.Time{}, err
	}
	// 2. 获取一个最新的区块 slot
	slot, err := s.rpcClient.GetSlot(context.TODO(), rpc.CommitmentFinalized)
	fmt.Printf("slot: %+v\n", slot)

	opts := &rpc.GetBlockOpts{
		// ！！！ 修复方案在这里 ！！！
		// 设置为 0 表示你支持 V0 版本交易，同时也兼容旧版 (Legacy) 交易
		MaxSupportedTransactionVersion: new(uint64), // 创建一个 *uint64
	}
	*opts.MaxSupportedTransactionVersion = 0 // 将其值设置为 0
	// block, err := s.rpcClient.GetBlock(context.TODO(), 415409862)
	block, err := s.rpcClient.GetBlockWithOpts(context.TODO(), 415409862, opts)

	for _, v := range block.Transactions {
		// var tx solana.Transaction
		// err := tx.UnmarshalBase64(v.Transaction.GetBinary())
		fmt.Printf("v: %+v\n", v.Transaction.GetBinary())
		// fmt.Printf("str: %+v\n", util.BytesToHexString(v.Transaction.GetBinary()))
		// fmt.Printf("str: %+v\n", util.Bytes2Hex(v.Transaction.GetBinary()))
	}
	fmt.Printf("err: %v\n", err)
	return int(blockhash.RPCContext.Context.Slot), time.Time{}, nil
}

func NewHTTPTransport(
	timeout time.Duration,
	maxIdleConnsPerHost int,
	keepAlive time.Duration,
) *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     timeout,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:1080",
		}),
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).Dial,
	}
}

func NewRPC(rpcEndpoint string) *rpc.Client {
	var (
		defaultMaxIdleConnsPerHost = 10
		defaultTimeout             = 25 * time.Second
		defaultKeepAlive           = 180 * time.Second
	)
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: NewHTTP(
			defaultTimeout,
			defaultMaxIdleConnsPerHost,
			defaultKeepAlive,
		),
	}
	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	return rpc.NewWithCustomRPCClient(rpcClient)
}

func NewHTTP(
	timeout time.Duration,
	maxIdleConnsPerHost int,
	keepAlive time.Duration,
) *http.Client {
	tr := NewHTTPTransport(
		timeout,
		maxIdleConnsPerHost,
		keepAlive,
	)

	return &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
}
