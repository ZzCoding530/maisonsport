package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitZap() {
	initLogger()
}

func initLogger() {
	// 配置日志级别
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // 设置日志级别为info

	// 配置日志输出
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 使用ISO8601时间格式
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建core，可以同时输出到多个地方
	infoFile := zapcore.AddSync(openLogFile("all.log"))
	errorFile := zapcore.AddSync(openLogFile("error.log"))
	stdout := zapcore.Lock(os.Stdout)
	stderr := zapcore.Lock(os.Stderr)

	// 根据配置创建logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoFile, zapcore.InfoLevel),
		zapcore.NewCore(encoder, errorFile, zapcore.ErrorLevel),
		zapcore.NewCore(encoder, stdout, zapcore.InfoLevel),
		zapcore.NewCore(encoder, stderr, zapcore.ErrorLevel),
	)

	Logger = zap.New(core)
}

func openLogFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
