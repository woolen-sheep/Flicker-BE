package log

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/woolen-sheep/Flicker-BE/config"
)

var Logger *logrus.Logger

func init() {
	Logger = getLogger()
	if Logger == nil {
		panic("Logger init failed")
	}
	Logger.Info("Logger init succeed")
}

func getLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.SetReportCaller(true)

	if config.C.Debug == true {
		logger.SetLevel(logrus.DebugLevel)
	}

	logConfig := config.C.LogConf

	baseLogPath := path.Join(logConfig.LogPath, logConfig.LogFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logger.Fatal(err)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	logger.AddHook(lfHook)
	return logger
}

func writeLog(fileName, funcName, errMsg, from string, err error) {
	Logger.WithFields(logrus.Fields{
		"package":  "package_name",
		"file":     fileName,
		"function": funcName,
		"err":      err,
		"from":     from,
	}).Warn(errMsg)
}
