package logger

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/constant"

	"go.uber.org/zap"
)

var Logger *zap.Logger
var ErrorLogger *zap.Logger
var TronLogger *zap.Logger
var TronShastaLogger *zap.Logger
var ArbitrumLogger *zap.Logger
var AvaxLogger *zap.Logger
var MaticLogger *zap.Logger
var BaseLogger *zap.Logger
var OpLogger *zap.Logger
var EthereumLogger *zap.Logger

func InitLogger(conf config.StorageConfig) {
	factory := NewLoggerFactory(conf)
	Logger = factory.GetLogger("default")
	ErrorLogger = factory.GetLogger("error")
	TronLogger = factory.GetLogger(constant.ChainTron)
	TronShastaLogger = factory.GetLogger(constant.ChainTronShasta)
	ArbitrumLogger = factory.GetLogger(constant.ChainArbitrum)
	AvaxLogger = factory.GetLogger(constant.ChainAvax)
	MaticLogger = factory.GetLogger(constant.ChainMatic)
	BaseLogger = factory.GetLogger(constant.ChainBase)
	OpLogger = factory.GetLogger(constant.ChainOp)
	EthereumLogger = factory.GetLogger(constant.ChainEthereum)
}
