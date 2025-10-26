package blockchain

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/util"
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
		Blockchain: *NewBlockchain(ctx, conf),
		rpcClient:  rpcClient,
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
	// block, err := s.rpcClient.GetBlockWithOpts()
	block, err := s.rpcClient.GetParsedBlockWithOpts(context.TODO(), 368735034, opts)
	for _, v := range block.Transactions {
		if len(v.Transaction.Signatures) == 0 {
			continue
		}
		if v.Transaction.Signatures[0].String() != "4YeDboDBZspmLa5fBDnzMQRg24koA5XD5vsCmutacvKnz2vCgE4qSHvYJCeU4xc1AmsdEAgyR7zh4txzuA7Rg9s2" {
			continue
		}
		// var tx solana.Transaction
		// err := tx.UnmarshalBase64(v.Transaction.GetBinary())
		// fmt.Printf("v.Transaction.Signatures: %+v\n", v.Transaction.Signatures)
		fmt.Printf("v: %+v\n", v.Transaction)
		// fmt.Printf("str: %+v\n", util.BytesToHexString(v.Transaction.GetBinary()))
		// fmt.Printf("str: %+v\n", util.Bytes2Hex(v.Transaction.GetBinary()))
	}
	fmt.Printf("err: %v\n", err)
	return int(blockhash.RPCContext.Context.Slot), time.Time{}, nil
}
func (s *Solana) GetTokenAccountBalance(address, contact_address, token_address string) (float64, string, error) {
	if token_address == "" {
		associatedTokenAddress, err := s.GetAssociatedTokenAddress(address, contact_address)
		if err != nil {
			return 0, "", err
		}
		token_address = associatedTokenAddress.String()
	}
	token_addr, _ := solana.PublicKeyFromBase58(token_address)
	balance, err := s.rpcClient.GetTokenAccountBalance(s.ctx, token_addr, rpc.CommitmentConfirmed)
	if err != nil {
		return 0, "", err
	}
	return *balance.Value.UiAmount, token_address, err
}
func (s *Solana) GetAssociatedTokenAddress(address string, contact_address string) (solana.PublicKey, error) {
	addr, _ := solana.PublicKeyFromBase58(address)
	contact_addr, _ := solana.PublicKeyFromBase58(contact_address)
	programAddress, _, err := solana.FindProgramAddress([][]byte{
		addr.Bytes(),
		solana.TokenProgramID.Bytes(),
		contact_addr.Bytes(),
	}, solana.SPLAssociatedTokenAccountProgramID)
	return programAddress, err
}

func (s *Solana) GetTokenChargeOfTransaction(signature string) {
	sign := solana.MustSignatureFromBase58(signature)
	tx, err := s.rpcClient.GetParsedTransaction(s.ctx, sign, &rpc.GetParsedTransactionOpts{
		Commitment:                     rpc.CommitmentConfirmed,
		MaxSupportedTransactionVersion: util.Uint64ToPtr(0),
	})
	pre_token_balances := tx.Meta.PreTokenBalances
	post_token_balances := tx.Meta.PostTokenBalances
	fmt.Printf("tx: %v\n", tx)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("pre_token_balances: %v\n", pre_token_balances)
	fmt.Printf("post_token_balances: %v\n", post_token_balances)
}

func (s *Solana) Test(token_address string) {
	signatures, err := s.rpcClient.GetSignaturesForAddress(s.ctx, solana.MustPublicKeyFromBase58(token_address))
	for _, v := range signatures {
		fmt.Printf("v: %+v\n", v)
		if v.Signature.String() != "9AdhVd9Ax9K32g8ErPgKRGh4YSedCpwtD86E3CWgx66Pq8kRx8iVqvj8muA25jTnVVLSYwgtHAaLav2AxsqCPjm" {
			continue
		}
		tx, err := s.rpcClient.GetParsedTransaction(s.ctx, v.Signature, &rpc.GetParsedTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: util.Uint64ToPtr(0),
		})

		fmt.Printf("tx: %+v\n", tx.Transaction)
		fmt.Printf("tx.Meta: %+v\n", tx.Meta)
		fmt.Printf("err: %v\n", err)
		break

	}
	// fmt.Printf("signatures: %+v\n", signatures)
	fmt.Printf("err: %v\n", err)
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
			// Host:   "127.0.0.1:1080",
			Host: "10.0.5.124:45613",
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
