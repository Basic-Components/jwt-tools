package logger

import (
	logrus "github.com/sirupsen/logrus"
)

// Init 初始化log的配置
func Init() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

//Logger 默认的logger
var Logger = Init()

//Log 有默认字段的log
var Log = Logger.WithFields(logrus.Fields{
	"app-type": "jwt-signer",
})
