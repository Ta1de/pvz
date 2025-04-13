package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - интерфейс логгера
type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Sync() error
}

type ZapSugaredLogger struct {
	Logger *zap.SugaredLogger // Делаем поле экспортируемым
}

var (
	Log Logger
)

// Init инициализирует глобальный логгер
func Init() error {
	logFile := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"

	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	sugared := zapLogger.Sugar()

	Log = &ZapSugaredLogger{Logger: sugared}
	zap.ReplaceGlobals(zapLogger)

	return nil
}

func (z *ZapSugaredLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.Logger.Infow(msg, keysAndValues...)
}

func (z *ZapSugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.Logger.Errorw(msg, keysAndValues...)
}

func (z *ZapSugaredLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.Logger.Fatalw(msg, keysAndValues...)
}

func (z *ZapSugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.Logger.Warnw(msg, keysAndValues...)
}

func (z *ZapSugaredLogger) Sync() error {
	return z.Logger.Sync()
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	formatted := t.Format("2006-01-02 15:04:05")
	enc.AppendString(formatted)
}
