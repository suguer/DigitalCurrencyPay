package blockchain

import (
	"DigitalCurrency/internal/blockchain/model"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/util"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

type Tron struct {
	Blockchain
	httpClient http.Client
	grpcClient *client.GrpcClient
}

func NewTron(ctx context.Context, conf *config.EthConfig) *Tron {
	tron := Tron{
		Blockchain{
			Chain:  conf.Name,
			Config: conf,
			ctx:    ctx,
		},
		http.Client{
			Timeout: 60 * time.Second, // 设置超时时间
		},
		client.NewGrpcClient(conf.GrpcAddress),
	}
	return &tron
}

func (t *Tron) GetNowBlockId(blockID ...int) (int, time.Time, error) {
	url := "/walletsolidity/getblock"
	method := "POST"
	params := make(map[string]any) // 请求参数
	params["detail"] = false
	if len(blockID) > 0 && blockID[0] > 0 { // 如果提供了区块 ID
		params["id_or_num"] = strconv.Itoa(blockID[0]) // 将 ID 转为字符串
	}
	res, err := t.request(method, url, params) // 发送请求
	if err != nil {
		return 0, time.Time{}, err // 返回错误
	}
	var data model.GetNowBlock
	err = json.Unmarshal(res, &data) // 解析响应
	if err != nil {
		return 0, time.Time{}, err // 返回错误
	}
	if data.BlockHeader.RawData.Number == 0 && data.Error != "" {
		return 0, time.Time{}, errors.New(data.Error) // 返回错误信息
	}
	return data.BlockHeader.RawData.Number, util.GetDateByTimestamp(data.BlockHeader.RawData.Timestamp), nil // 返回区块 ID 和时间
}

func (t *Tron) GetBlockByNum(num int) (model.GetBlockByNum, error) {

	url := "/walletsolidity/getblockbynum"
	params := make(map[string]any) // 请求参数
	params["num"] = num
	res, err := t.request("POST", url, params) // 发送请求
	var data model.GetBlockByNum
	if err != nil {
		return data, err // 返回错误
	}
	err = json.Unmarshal(res, &data) // 解析响应
	return data, err                 // 返回响应体和错误信息
}

func (t *Tron) GetAccountResources(address string) (data model.AccountResource, err error) {
	if !strings.HasPrefix(address, "0x41") {
		address = "0x41" + address[2:]
	}
	url := "/wallet/getaccountresource"
	params := make(map[string]any) // 请求参数
	params["address"] = address
	res, err := t.request("POST", url, params) // 发送请求
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res, &data) // 解析响应
	if err != nil {
		return data, err
	}
	if data.TotalEnergyLimit == 0 {
		return data, errors.New("account not exist")
	}
	return data, nil
}

func (t *Tron) signTransaction(rawData []byte, privateKeyHex string) ([]byte, error) {
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	sk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(hash, sk)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
func (t *Tron) broadcastTransaction(tx *api.TransactionExtention) (string, error) {
	result, err := t.grpcClient.Broadcast(tx.Transaction)
	if err != nil {
		return "", err
	}
	if !result.Result {
		return "", errors.New("broadcast failed")
	}
	txid := tx.GetTxid()
	return util.BytesToHexString(txid), nil
}

func (t *Tron) TransferTrx(privateKeyHex, ownerAddress, toAddress string, amount int64) (txId string, err error) {
	err = t.grpcClient.Start(client.GRPCInsecure())
	if err != nil {
		return "", err
	}
	ownerAddr, err := util.HexString2Address(ownerAddress)
	if err != nil {
		ownerAddr = ownerAddress
	}
	toAddr, err := util.HexString2Address(toAddress)
	if err != nil {
		toAddr = toAddress
	}
	tx, err := t.grpcClient.Transfer(ownerAddr, toAddr, amount)
	if err != nil {
		return "", err
	}
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return "", err
	}
	signature, err := t.signTransaction(rawData, privateKeyHex)
	if err != nil {
		return "", err
	}
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)
	return t.broadcastTransaction(tx)
}

func (t *Tron) TransferTRC20(privateKeyHex, ownerAddress, contractAddress, toAddress string, amount int64) (txId string, err error) {
	err = t.grpcClient.Start(client.GRPCInsecure())
	if err != nil {
		return "", err
	}
	amount = amount * int64(math.Pow10(t.Blockchain.Config.Precision))
	tx, err := t.grpcClient.TRC20Send(ownerAddress, toAddress, contractAddress, big.NewInt(amount), 50)
	if err != nil {
		return "", err
	}
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return "", err
	}
	signature, err := t.signTransaction(rawData, privateKeyHex)
	if err != nil {
		return "", err
	}
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)
	return t.broadcastTransaction(tx)
}

func (t *Tron) request(method string, action string, param any) ([]byte, error) {
	jsonData, err := json.Marshal(param) // 将参数序列化为 JSON
	if err != nil {
		return nil, err // 返回错误
	}
	var reqBody io.Reader
	uri := fmt.Sprintf("%v%v", t.Config.Node, action) // 构建完整 URL
	if method == "GET" {
		reqBody = nil
		params := url.Values{}
		if param != nil {
			for k, v := range param.(map[string]any) {
				params.Add(k, fmt.Sprintf("%v", v))
			}
		}
		uri = fmt.Sprintf("%v%v?%v", t.Config.Node, action, params.Encode())
	} else {
		reqBody = bytes.NewBuffer(jsonData) // 创建请求体
	}
	req, err := http.NewRequest(method, uri, reqBody) // 创建 HTTP 请求
	if len(t.Config.Key) > 0 {
		req.Header.Add("TRON-PRO-API-KEY", t.Config.Key) // 添加 API 密钥
	}
	if err != nil {
		return nil, err // 返回错误
	}
	httpRes, err := t.httpClient.Do(req) // 发送请求
	if err != nil {
		return nil, err // 返回错误
	}
	defer httpRes.Body.Close()            // 确保在函数结束时关闭响应体
	body, err := io.ReadAll(httpRes.Body) // 读取响应体
	return body, err                      // 返回响应体和错误信息
}

func (t *Tron) requestV2(method string, action string, reqByte []byte) ([]byte, error) {
	reqBody := bytes.NewBuffer(reqByte)
	uri := fmt.Sprintf("%v%v", t.Config.Node, action) // 构建完整 URL
	req, err := http.NewRequest(method, uri, reqBody) // 创建 HTTP 请求
	if len(t.Config.Key) > 0 {
		req.Header.Add("TRON-PRO-API-KEY", t.Config.Key) // 添加 API 密钥
	}
	if err != nil {
		return nil, err // 返回错误
	}
	httpRes, err := t.httpClient.Do(req) // 发送请求
	if err != nil {
		return nil, err // 返回错误
	}
	defer httpRes.Body.Close()            // 确保在函数结束时关闭响应体
	body, err := io.ReadAll(httpRes.Body) // 读取响应体
	return body, err                      // 返回响应体和错误信息
}

func (t *Tron) TokenTx(toAddress, contractAddress string) (model.TokenTxResponse, error) {
	toAddr, _ := util.HexString2Address(toAddress)
	contractAddr, _ := util.HexString2Address(contractAddress)

	url := fmt.Sprintf("/v1/accounts/%v/transactions/trc20", toAddr)
	method := "GET"
	params := make(map[string]any) // 请求参数
	params["only_confirmed"] = 1
	params["limit"] = 100
	params["contract_address"] = contractAddr
	params["order_by"] = "block_timestamp,desc"
	params["only_to"] = 1
	res, err := t.request(method, url, params) // 发送请求
	var tronTx model.TronTokenTxResponse
	var tokenTx model.TokenTxResponse
	if err != nil {
		return tokenTx, err // 返回错误
	}
	err = json.Unmarshal(res, &tronTx) // 解析响应
	if err != nil {
		return model.TokenTxResponse{}, err // 返回错误
	}
	tokenTx.Result = make([]model.TokenTxTransaction, len(tronTx.Data))
	for i, tx := range tronTx.Data {

		fromAddr, _ := address.Base58ToAddress(tx.From)
		toAddr, _ := address.Base58ToAddress(tx.To)
		contractAddr, _ := address.Base58ToAddress(tx.TokenInfo.Address)
		tokenTx.Result[i].TimeStamp = fmt.Sprintf("%d", tx.BlockTimestamp/1000)
		tokenTx.Result[i].Hash = tx.TransactionId
		tokenTx.Result[i].From = "0x" + fromAddr.Hex()[4:]
		tokenTx.Result[i].To = "0x" + toAddr.Hex()[4:]
		tokenTx.Result[i].Value = tx.Value
		tokenTx.Result[i].TokenDecimal = strconv.FormatInt(tx.TokenInfo.Decimals, 10)
		tokenTx.Result[i].TokenName = tx.TokenInfo.Name
		tokenTx.Result[i].TokenSymbol = tx.TokenInfo.Symbol
		tokenTx.Result[i].ContractAddress = contractAddr.Hex()
	}
	return tokenTx, nil
}
