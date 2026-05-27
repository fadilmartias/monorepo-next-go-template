package logger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Internal logger instance
var log *zap.SugaredLogger

// Export writer untuk HTTP middleware Fiber
var FileWriter *lumberjack.Logger

func Init() {
	os.MkdirAll("./storage/logs", 0755)

	FileWriter = &lumberjack.Logger{
		Filename:   "./storage/logs/app.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeTime:     zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleCfg := encoderConfig
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	fileCfg := encoderConfig
	fileCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileCfg)

	// --- OPTIMASI ASYNC DI SINI ---
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Bungkus FileWriter dengan BufferedWriteSyncer
	asyncFileWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(FileWriter),
		Size:          256 * 1024,      // Buffer sebesar 256 KB di memori
		FlushInterval: 2 * time.Second, // Tulis ke disk otomatis setiap 2 detik
	}

	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)

	// Gunakan asyncFileWriter
	fileCore := zapcore.NewCore(fileEncoder, asyncFileWriter, zapcore.ErrorLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	baseLogger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	log = baseLogger.Sugar()
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// --- FUNGSI UNTUK MENGAMBIL REQUEST ID DARI FIBER ---
// Gunakan ini di dalam handler/controller
func Ctx(c fiber.Ctx) *zap.SugaredLogger {
	// Fiber requestid middleware biasanya menyimpan ID di Locals("requestid")
	reqID := requestid.FromContext(c)
	if reqID != "" {
		return log.With("req_id", reqID)
	}
	return log
}

func Base() *zap.SugaredLogger {
	return log
}

// --- FUNGSI GLOBAL (TANPA REQUEST ID) ---
func Debug(args ...any)                   { log.Debug(args...) }
func Debugf(template string, args ...any) { log.Debugf(template, args...) }

func Info(args ...any)                   { log.Info(args...) }
func Infof(template string, args ...any) { log.Infof(template, args...) }

func Warn(args ...any)                   { log.Warn(args...) }
func Warnf(template string, args ...any) { log.Warnf(template, args...) }

func Error(args ...any)                   { log.Error(args...) }
func Errorf(template string, args ...any) { log.Errorf(template, args...) }
