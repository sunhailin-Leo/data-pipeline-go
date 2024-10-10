package logger

import (
	"io"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// Logger logger object
var Logger *zap.Logger

// logger time encoder function
var timeEncoderFunc = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// logger duration encoder function
var durationEncoderFunc = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(int64(d) / 1000000)
}

// logger json encoder function
var jsonEncoderFunc = func(writer io.Writer) zapcore.ReflectedEncoder {
	enc := sonic.ConfigDefault.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	return enc
}

// logger function caller encoder function
var callerEncoderFunc = func(isAddServiceName bool) zapcore.CallerEncoder {
	return func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		if isAddServiceName {
			encoder.AppendString(utils.ServiceName)
		}
		encoder.AppendString("\"" + caller.String() + "\"")
	}
}

// geo zap EncoderConfig
func getEncodeConfig(isAddServiceName bool) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:             "time",
		LevelKey:            "level",
		NameKey:             "logger",
		CallerKey:           "linenum",
		FunctionKey:         "function",
		MessageKey:          "msg",
		StacktraceKey:       "stacktrace",
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          timeEncoderFunc,
		EncodeDuration:      durationEncoderFunc,
		EncodeCaller:        callerEncoderFunc(isAddServiceName),
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: jsonEncoderFunc,
		ConsoleSeparator:    utils.LogConsoleSeparator,
	}
}

// get logger rolling config
func getRollingConfig(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    utils.LogMaxSize,
		MaxAge:     utils.LogMaxAge,
		MaxBackups: utils.LogMaxBackups,
		Compress:   true,
	}
}

// NewZapLogger init zap Logger
func NewZapLogger() {
	Logger = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(getEncodeConfig(false)), // encoder config
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(getRollingConfig("logs/"+utils.ServiceName+"/"+utils.GetHostName()+".log")),
			), // writer to console and file
			zap.NewAtomicLevelAt(zap.DebugLevel), // set logger level
		),
		zap.AddCaller(),   // enable stack trace
		zap.Development(), // enable development mode
	)
}
