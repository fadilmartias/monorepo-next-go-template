package logger

import (
	"os"

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
		MaxSize:    10, // 10 MB
		MaxBackups: 30,
		MaxAge:     30, // 30 hari
		Compress:   true,
		LocalTime:  true,
	}

	// 1. Konfigurasi Dasar Encoder (Format "Mudah Dibaca")
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeTime:     zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 2. Setup Console Encoder (DENGAN WARNA)
	consoleCfg := encoderConfig
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Warna aktif!
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	// 3. Setup File Encoder (TANPA WARNA - agar file bersih)
	fileCfg := encoderConfig
	fileCfg.EncodeLevel = zapcore.CapitalLevelEncoder // Tanpa warna!
	fileEncoder := zapcore.NewConsoleEncoder(fileCfg) // Gunakan format Console (bukan JSON) agar mudah dibaca manusia

	// Writers
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileZapWriter := zapcore.AddSync(FileWriter)

	// 4. Aturan Level Log
	// Console: Tampilkan semua log dari level Debug hingga Error
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)

	// File: TAMPILKAN INFO, WARN, dan ERROR.
	// (Kita pakai InfoLevel agar kamu bisa pilih manual log mana yang mau dimasukkan ke file)
	fileCore := zapcore.NewCore(fileEncoder, fileZapWriter, zapcore.InfoLevel)

	// Gabungkan
	core := zapcore.NewTee(consoleCore, fileCore)

	// AddStacktrace hanya akan print detail error saat benar-benar terjadi Error
	baseLogger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	log = baseLogger.Sugar()
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

// --- FUNGSI GLOBAL (TANPA REQUEST ID) ---
func Debug(args ...any)                   { log.Debug(args...) }
func Debugf(template string, args ...any) { log.Debugf(template, args...) }

func Info(args ...any)                   { log.Info(args...) }
func Infof(template string, args ...any) { log.Infof(template, args...) }

func Warn(args ...any)                   { log.Warn(args...) }
func Warnf(template string, args ...any) { log.Warnf(template, args...) }

func Error(args ...any)                   { log.Error(args...) }
func Errorf(template string, args ...any) { log.Errorf(template, args...) }
