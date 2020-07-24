package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
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
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{file, os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	}
}
func Trace(args ...interface{}) {
	logrus.Trace(args...)
}
func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}
func Traceln(args ...interface{}) {
	logrus.Traceln(args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}
func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}
func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
func Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}
