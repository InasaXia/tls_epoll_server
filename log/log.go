package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

var Logger *zap.SugaredLogger

func getWriter(file string) io.Writer {
	writer,err := rotatelogs.New(
		strings.Replace(file,".log","",1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(file),
		rotatelogs.WithMaxAge(time.Hour*12*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err!=nil {
		panic(err)
	}
	return writer
}
func InitLogger(path string){
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		//NameKey:        "",
		CallerKey:      "file",
		//StacktraceKey:  "",
		//LineEnding:     "",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeName:     nil,
	})
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})
	infoWriter := getWriter(path+"/info.log")
	errorWriter := getWriter(path+"/error.log")
	core := zapcore.NewTee(
		zapcore.NewCore(encoder,zapcore.AddSync(infoWriter),infoLevel),
		zapcore.NewCore(encoder,zapcore.AddSync(errorWriter),errorLevel),
	)
	log := zap.New(core,zap.AddCaller())
	Logger = log.Sugar()
}


