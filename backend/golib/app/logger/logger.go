package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"learnerai/golib/app/context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// error logger
var errorLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

const ErrorDir = "error/"

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func getWriter(filename string, expireDay int32, format string) io.Writer {
	if format == "" {
		format = "%Y%m%d%H"
	}
	hook, err := rotatelogs.New(
		filename+"."+format,
		rotatelogs.WithMaxAge(time.Duration(expireDay)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

// Setup initialize the log instance
func Setup(filePath, fileName, logLevel, runMode string, expireDay int32) {
	fileFullName := filePath + fileName
	errFullName := filePath + ErrorDir + fileName + "_error"
	level := getLoggerLevel(logLevel)
	debugLevel := zap.DebugLevel

	//日志切割配置
	syncWriter := getWriter(fileFullName, expireDay, "")
	syncErrorWriter := getWriter(errFullName, expireDay, "")

	//日志编码配置
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//真正的来配置zap
	core := zapcore.NewTee(
		zapcore.NewCore( //文件，json
			zapcore.NewJSONEncoder(encoder),
			zapcore.AddSync(syncWriter),
			zap.NewAtomicLevelAt(level),
		),
		zapcore.NewCore( // 控制台
			zapcore.NewConsoleEncoder(encoder),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(debugLevel),
		),
		zapcore.NewCore( //文件，json
			zapcore.NewJSONEncoder(encoder),
			zapcore.AddSync(syncErrorWriter),
			zap.NewAtomicLevelAt(zapcore.WarnLevel),
		),
	)
	var logger *zap.Logger
	additionalFields := zap.Fields(
		zap.Int("pid", os.Getpid()),
		zap.String("process", path.Base(os.Args[0])),
	)
	if runMode == "debug" {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development(), additionalFields)
	} else {
		logger = zap.New(core)
	}
	//logger := zap.New(core, zap.AddCaller(), zap.Development())
	errorLogger = logger.Sugar()
}

// 增加文件切割规则配置
func Setup2(filePath, fileName, logLevel, runMode string, expireDay int32, format string) {
	fileFullName := filePath + fileName
	level := getLoggerLevel(logLevel)
	debugLevel := zap.DebugLevel

	//日志切割配置
	syncWriter := getWriter(fileFullName, expireDay, format)

	//日志编码配置
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//真正的来配置zap
	core := zapcore.NewTee(
		zapcore.NewCore( //文件，json
			zapcore.NewJSONEncoder(encoder),
			zapcore.AddSync(syncWriter),
			zap.NewAtomicLevelAt(level),
		),
		zapcore.NewCore( // 控制台
			zapcore.NewConsoleEncoder(encoder),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(debugLevel),
		),
	)
	var logger *zap.Logger
	additionalFields := zap.Fields(
		zap.Int("pid", os.Getpid()),
		zap.String("process", path.Base(os.Args[0])),
	)
	if runMode == "debug" {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development(), additionalFields)
	} else {
		logger = zap.New(core)
	}
	//logger := zap.New(core, zap.AddCaller(), zap.Development())
	errorLogger = logger.Sugar()
}

// 日志不输出到屏幕的初始化配置
func SetupWithNoPrintln(filePath, fileName, logLevel, runMode string, expireDay int32) {
	fileFullName := filePath + fileName
	level := getLoggerLevel(logLevel)
	//debugLevel := zap.DebugLevel

	//日志切割配置
	syncWriter := getWriter(fileFullName, expireDay, "")

	//日志编码配置
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//真正的来配置zap
	core := zapcore.NewTee(zapcore.NewCore( //文件，json
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(syncWriter),
		zap.NewAtomicLevelAt(level),
	))
	var logger *zap.Logger
	additionalFields := zap.Fields(
		zap.Int("pid", os.Getpid()),
		zap.String("process", path.Base(os.Args[0])),
	)
	if runMode == "debug" {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development(), additionalFields)
	} else {
		logger = zap.New(core)
	}
	//logger := zap.New(core, zap.AddCaller(), zap.Development())
	errorLogger = logger.Sugar()
}

func getCommonLogString(ctx *gin.Context, template string) string {
	var globalLogid, url, remoteAddress string
	if ctx != nil {
		globalLogid = context.GetRequestID(ctx)
		if ctx.Request != nil {
			url = ctx.Request.RequestURI
			remoteAddress = ctx.Request.RemoteAddr
		}
	}

	return fmt.Sprintf("[logid=%s] [uri=%s] [remote_ip=%s] %s", globalLogid, url, remoteAddress, template)
}

func getJsonArgs(args ...interface{}) []interface{} {

	newArgs := []interface{}{}
	for _, arg := range args {
		switch t := arg.(type) {
		case error:
			newArgs = append(newArgs, arg.(error).Error())
		default:
			j, err := json.Marshal(arg)
			if err != nil {
				fmt.Printf("json marshal error:%+v, t:%+v", err, t)
			}
			newArgs = append(newArgs, string(j))
		}
	}

	return newArgs
}

func Debug(ctx *gin.Context, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	commonLogString := getCommonLogString(ctx, "")
	newArgs := []interface{}{commonLogString}
	newArgs = append(newArgs, args...)
	errorLogger.Debug(newArgs...)
}

func Debugf(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Debugf(newTemplate, args...)
}

func DebugfJson(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newArgs := getJsonArgs()
	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Debugf(newTemplate, newArgs...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	commonLogString := getCommonLogString(ctx, "")
	newArgs := []interface{}{commonLogString}
	newArgs = append(newArgs, args...)
	errorLogger.Info(newArgs...)
}

func Infof(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Infof(newTemplate, args...)
}

func InfofJson(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newArgs := getJsonArgs()
	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Infof(newTemplate, newArgs...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	commonLogString := getCommonLogString(ctx, "")
	newArgs := []interface{}{commonLogString}
	newArgs = append(newArgs, args...)
	errorLogger.Warn(newArgs...)
}

func Warnf(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Warnf(newTemplate, args...)
}

func WarnfJson(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newArgs := getJsonArgs()
	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Warnf(newTemplate, newArgs...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	commonLogString := getCommonLogString(ctx, "")
	newArgs := []interface{}{commonLogString}
	newArgs = append(newArgs, args...)
	errorLogger.Error(newArgs...)
}

func Errorf(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Errorf(newTemplate, args...)
}

func ErrorfJson(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newArgs := getJsonArgs()
	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Errorf(newTemplate, newArgs...)
}

// fatal不能使用，进程会直接退出
func Fatal(ctx *gin.Context, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	commonLogString := getCommonLogString(ctx, "")
	newArgs := []interface{}{commonLogString}
	newArgs = append(newArgs, args...)
	errorLogger.Error(newArgs...)
}

func Fatalf(ctx *gin.Context, template string, args ...interface{}) {
	if errorLogger == nil {
		return
	}

	newTemplate := getCommonLogString(ctx, template)
	errorLogger.Errorf(newTemplate, args...)
}
