package log

import (
	"github.com/iok200/go-ok/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type Config struct {
	Name string `properties:"log.name,default=info.log"`
}

var _log_is_init = false
var _log_mu sync.Mutex

func initLog() {
	if _log_is_init {
		return
	}
	_log_mu.Lock()
	defer _log_mu.Unlock()
	if _log_is_init {
		return
	}
	_log_is_init = true
	var cfg Config
	config.Parse(&cfg)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             true,
		ForceQuote:                false,
		DisableQuote:              true,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05.000",
		DisableSorting:            true,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	file, err := os.OpenFile(cfg.Name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{file, os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	}
}

func Trace(args ...interface{}) {
	initLog()
	logrus.Trace(args...)
}
func Tracef(format string, args ...interface{}) {
	initLog()
	logrus.Tracef(format, args...)
}
func Traceln(args ...interface{}) {
	initLog()
	logrus.Traceln(args...)
}

func Debug(args ...interface{}) {
	initLog()
	logrus.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	initLog()
	logrus.Debugf(format, args...)
}
func Debugln(args ...interface{}) {
	initLog()
	logrus.Debugln(args...)
}

func Info(args ...interface{}) {
	initLog()
	logrus.Info(args...)
}
func Infof(format string, args ...interface{}) {
	initLog()
	logrus.Infof(format, args...)
}
func Infoln(args ...interface{}) {
	initLog()
	logrus.Infoln(args...)
}

func Warn(args ...interface{}) {
	initLog()
	logrus.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	initLog()
	logrus.Warnf(format, args...)
}
func Warnln(args ...interface{}) {
	initLog()
	logrus.Warnln(args...)
}

func Error(args ...interface{}) {
	initLog()
	logrus.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	initLog()
	logrus.Errorf(format, args...)
}
func Errorln(args ...interface{}) {
	initLog()
	logrus.Errorln(args...)
}

func Fatal(args ...interface{}) {
	initLog()
	logrus.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	initLog()
	logrus.Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
	initLog()
	logrus.Fatalln(args...)
}

func Panic(args ...interface{}) {
	initLog()
	logrus.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	initLog()
	logrus.Panicf(format, args...)
}
func Panicln(args ...interface{}) {
	initLog()
	logrus.Panicln(args...)
}
