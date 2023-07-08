package logger

import (
	"bcloud/netdisk/floder"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger(level, mode string) {
	writeSyncer := getLogWriter(fmt.Sprintf("%s/bcloud.log", floder.GetConfigDir()), 200, 30, 7)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(level))
	if err != nil {
		return
	}
	var core zapcore.Core
	if mode == "dev" {
		//	开发模式 日志输出到终端和文件中
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			// 到终端
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			// 到文件
			zapcore.NewCore(encoder, writeSyncer, l),
		)
	} else {
		// 日志收入到文件
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
