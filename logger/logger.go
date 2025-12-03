package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers   = make(map[string]*logrus.Logger)
	once      sync.Once
	basePath  string
	initMutex sync.Mutex
)

// ---------------- Parse Log Level --------------------
func parseLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// ---------------- Init Base Path --------------------
func initBasePath() {
	once.Do(func() {
		basePath = filepath.Join("logs")
		_ = os.MkdirAll(basePath, 0755)
	})
}

// ---------------- Create Logger Instance --------------------
func getLogger(folder string) *logrus.Logger {
	initBasePath()

	initMutex.Lock()
	defer initMutex.Unlock()

	if logger, exists := loggers[folder]; exists {
		return logger
	}

	logDir := filepath.Join(basePath, folder)
	_ = os.MkdirAll(logDir, 0755)

	// Format file name theo ngày
	today := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("app_%s.log", today)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, logFileName),
		MaxSize:    10, // MB before rotating
		MaxBackups: 7,  // keep 7 rotated files
		MaxAge:     60, // days old
		Compress:   true,
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output both console & file
	logger.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
	logger.SetLevel(logrus.DebugLevel)

	loggers[folder] = logger
	return logger
}

// ---------------- Public Functions --------------------

func WriteLogMsg(level string, msg string) {
	getLogger("LogMsg").Log(parseLogLevel(level), msg)
}

func WriteLogData(level string, data any) {
	getLogger("LogData").WithField("data", data).Log(parseLogLevel(level), "Data log")
}

func WriteLogEx(level string, msg string, data any) {
	getLogger("LogEx").WithField("data", data).Log(parseLogLevel(level), msg)
}
