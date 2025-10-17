package logger

import (
	"DigitalCurrency/internal/config"
	"fmt"
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
)

type LoggerFactory struct {
	Conf config.StorageConfig
}

func NewLoggerFactory(conf config.StorageConfig) *LoggerFactory {
	var factory *LoggerFactory
	once.Do(func() {
		factory = &LoggerFactory{Conf: conf}
	})
	return factory
}

func (f *LoggerFactory) GetLogger(name string) *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 添加以下配置来显示调用者信息
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 显示调用者信息
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 创建一个日志核心，输出到文件
	fileCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", f.Conf.Log, name),
			MaxSize:    30, // MB
			MaxBackups: 3,
			MaxAge:     7, // days
		}),
		zap.InfoLevel,
	)

	// 创建另一个日志核心，输出到标准输出
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.ErrorLevel,
	)

	// 使用 zapcore.NewTee 将两个核心组合起来
	logger := zap.New(zapcore.NewTee(fileCore, consoleCore),
		zap.AddCaller(),      // 添加调用者信息
		zap.AddCallerSkip(0), // 调整调用栈跳过的帧数
	)
	return logger
}
