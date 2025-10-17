package blockchain

import (
	"DigitalCurrency/internal/blockchain/model"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/constant"
	"DigitalCurrency/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type BlockchainInterface interface {
	TokenTx(toAddress, contractAddress string) (model.TokenTxResponse, error)
}

type Blockchain struct {
	Chain      string
	Config     *config.EthConfig
	ctx        context.Context
	client     *ethclient.Client
	httpClient *http.Client
}

func NewBlockchain(ctx context.Context, conf *config.EthConfig) *Blockchain {
	client, _ := ethclient.Dial(conf.GrpcAddress)
	httpClient := &http.Client{
		Timeout: time.Second * 30,
		// 设置代理
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "10.0.5.124:45613",
				// Host:   "127.0.0.1:1080",
			}),
		},
	}
	return &Blockchain{
		Chain:      conf.Name,
		Config:     conf,
		ctx:        ctx,
		client:     client,
		httpClient: httpClient,
	}

}
func (t *Blockchain) GetNowBlockId(blockID ...int) (int, time.Time, error) {
	num, err := t.client.BlockNumber(t.ctx)
	return int(num), time.Time{}, err
}

func (t *Blockchain) GetBlockByNum(num int) {
	block, err := t.client.BlockByNumber(t.ctx, big.NewInt(int64(num)))
	fmt.Printf("block: %v\n", block)
	fmt.Printf("err: %v\n", err)
}

func (t *Blockchain) GetLogs(FromBlock, ToBlock int) ([]types.Log, error) {
	logs, err := t.client.FilterLogs(t.ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(FromBlock)),
		ToBlock:   big.NewInt(int64(ToBlock)),
		Topics: [][]common.Hash{
			{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
		},
	})
	return logs, err
}

func (t *Blockchain) GetLogsFormat(FromBlock, ToBlock int) ([]model.EthLog, error) {
	logs, err := t.GetLogs(FromBlock, ToBlock)
	if err != nil {
		return nil, err
	}
	var ethLogs []model.EthLog
	for _, log := range logs {
		if len(log.Topics) != 3 {
			continue
		}
		// 创建一个新的 big.Int 实例
		bigInt := new(big.Int)
		bigInt.SetString(util.BytesToHexString(log.Data), 16)
		dec, _ := strconv.ParseUint(util.BytesToHexString(log.Data), 16, 64)
		ethLogs = append(ethLogs, model.EthLog{
			TxHash:          log.TxHash.String(),
			Amount:          uint(dec),
			FromAddress:     t.formatEthereumAddress(log.Topics[1].Hex()),
			ToAddress:       t.formatEthereumAddress(log.Topics[2].Hex()),
			ContractAddress: log.Address.Hex(),
			BlockNumber:     log.BlockNumber,
			BlockHash:       log.BlockHash.String(),
		})
	}
	return ethLogs, nil
}

func (t *Blockchain) EthGasFeeEstimate(fromAddress, toAddress, contractAddress string, amount uint64) (model.EthGasFee, error) {
	lastBlock, err := t.client.BlockByNumber(t.ctx, big.NewInt(int64(rpc.SafeBlockNumber)))
	fmt.Printf("lastBlock: %v\n", lastBlock)
	fmt.Printf("err: %v\n", err)
	return model.EthGasFee{}, nil
}
func (t *Blockchain) EthTransactionSend(fromAddress, toAddress, contractAddress string, amount uint64) error {

	toAddr := common.HexToAddress(toAddress)
	intAmount := new(big.Int).SetUint64(amount)
	data := util.EthContractTransferDataEncode(toAddr.Hex(), intAmount)
	nonce, _ := t.GetTransactionCount(toAddress)
	estimate, err := t.EthGasFeeEstimate(fromAddress, toAddress, contractAddress, amount)
	transaction := types.NewTransaction(
		nonce,
		toAddr,
		intAmount,
		estimate.GasLimit,
		new(big.Int).SetUint64(0),
		common.FromHex(data),
	)
	err = t.client.SendTransaction(t.ctx, transaction)
	hash := transaction.Hash()
	fmt.Printf("hash: %v\n", hash)
	fmt.Printf("err: %v\n", err)
	return err

}

func (t *Blockchain) GetTransactionCount(address string) (uint64, error) {
	// get_transaction_count
	nonce, err := t.client.NonceAt(t.ctx, common.HexToAddress(address), nil)
	return nonce, err
}

func (t *Blockchain) EthTransactionReceiptGet(tx string) (*types.Receipt, error) {
	tx_hash := common.HexToHash(tx)
	data, err := t.client.TransactionReceipt(t.ctx, tx_hash)
	return data, err
}

func (t *Blockchain) TokenTx(toAddress, contractAddress string) (model.TokenTxResponse, error) {
	var apiResponse model.TokenTxResponse
	params := url.Values{}
	params.Add("chainid", strconv.Itoa(t.Config.ChainId))
	params.Add("module", "account")
	params.Add("action", "tokentx")
	params.Add("page", "1")
	params.Add("offset", "1")
	params.Add("sort", "desc")
	params.Add("apikey", t.Config.Key)
	params.Add("address", toAddress)
	params.Add("contractAddress", contractAddress)
	URL := t.Config.Node + "?" + params.Encode()
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return apiResponse, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return apiResponse, err
	}
	defer resp.Body.Close()
	// 读取返回结果
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse, err
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return apiResponse, err
	}

	return apiResponse, nil
}

func Factory(ctx context.Context, chain string) BlockchainInterface {
	switch chain {
	case constant.ChainTron:
		return NewTron(ctx, &config.Conf.BlockChain.Tron)
	case constant.ChainTronShasta:
		return NewTron(ctx, &config.Conf.BlockChain.TronShasta)
	default:
		return NewBlockchain(ctx, &config.Conf.BlockChain.Arbitrum)
	}
}

func (t *Blockchain) formatEthereumAddress(hexStr string) string {
	// 移除 "0x" 前缀
	cleaned := strings.TrimPrefix(hexStr, "0x")

	// 找到第一个非零字符的位置
	firstNonZero := -1
	for i, char := range cleaned {
		if char != '0' {
			firstNonZero = i
			break
		}
	}

	// 如果全是零，返回零地址
	if firstNonZero == -1 {
		return "0x0"
	}

	// 提取非零部分
	nonZeroPart := cleaned[firstNonZero:]

	// 确保长度为40个字符（以太坊地址长度）
	if len(nonZeroPart) < 40 {
		// 如果长度不足，在前面补零
		nonZeroPart = strings.Repeat("0", 40-len(nonZeroPart)) + nonZeroPart
	} else if len(nonZeroPart) > 40 {
		// 如果长度超过，截取后40位
		nonZeroPart = nonZeroPart[len(nonZeroPart)-40:]
	}

	// 添加 "0x" 前缀
	return "0x" + nonZeroPart
}
