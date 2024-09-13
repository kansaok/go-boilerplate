package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kansaok/go-boilerplate/pkg/metadata"
	"github.com/sirupsen/logrus"
	"github.com/vardius/golog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Default values in case Caller is nil
	file := "unknown file"
	line := 0

	if entry.Caller != nil {
		file = entry.Caller.File
		line = entry.Caller.Line
	}

	return []byte(fmt.Sprintf(
		"ERROR - %s \nMessage: %s\nFile: %s\nLine: %d\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Message,
		file,
		line,
	)), nil
}

var (
	ErrorLogger *logrus.Logger
	InfoLogger  *logrus.Logger
	Logger      golog.Logger
)

func SetFlags(flag int) {
	Logger.SetFlags(flag)
}

func SetVerbosity(verbosity golog.Verbose) {
	Logger.SetVerbosity(verbosity)
}

func Debug(ctx context.Context, v string) {
	Logger.Debug(ctx, messageWithMeta(ctx, v))
}

func Info(ctx context.Context, v string) {
	Logger.Info(ctx, messageWithMeta(ctx, v))
}

func Warning(ctx context.Context, v string) {
	Logger.Warning(ctx, messageWithMeta(ctx, v))
}

func Error(ctx context.Context, v string) {
	Logger.Error(ctx, messageWithMeta(ctx, v))
}

func Critical(ctx context.Context, v string) {
	Logger.Critical(ctx, messageWithMeta(ctx, v))
}

func Fatal(ctx context.Context, v string) {
	Logger.Fatal(ctx, messageWithMeta(ctx, v))
}

func messageWithMeta(ctx context.Context, v string) string {
	mtd, _ := metadata.FromContext(ctx)
	s, err := json.Marshal(struct {
		Message string             `json:"message"`
		Meta    *metadata.Metadata `json:"meta"`
	}{
		Message: v,
		Meta:    mtd,
	})
	if err == nil {
		v = string(s)
	}

	return v
}

func Init() {
	// Create the log directory if it does not exist
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		panic(err)
	}

	l := golog.New()
	l.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)

	Logger = l

	// Initialize Error Logger with lumberjack for log rotation
	ErrorLogger = logrus.New()
	ErrorLogger.SetLevel(logrus.ErrorLevel)
	ErrorLogger.SetFormatter(&CustomFormatter{})
	errorLogWriter := &lumberjack.Logger{
		Filename:   "log/error.log",
		MaxSize:    10,  // Maximum size in MB before rotation
		MaxBackups: 3,   // Maximum number of backups to keep
		MaxAge:     30,  // Maximum age in days
		Compress:   true, // Compress log files
	}
	ErrorLogger.SetOutput(errorLogWriter)

	// Initialize Info Logger with lumberjack for log rotation
	InfoLogger = logrus.New()
	InfoLogger.SetLevel(logrus.InfoLevel)
	InfoLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
	})
	infoLogWriter := &lumberjack.Logger{
		Filename:   "log/info.log",
		MaxSize:    10,  // Maximum size in MB before rotation
		MaxBackups: 3,   // Maximum number of backups to keep
		MaxAge:     30,  // Maximum age in days
		Compress:   true, // Compress log files
	}
	InfoLogger.SetOutput(infoLogWriter)
}
