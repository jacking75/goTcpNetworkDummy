package utils

import (
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"encoding/json"

	"go.uber.org/zap"
)

var (
	Logger, _ = zap.NewProduction()
)

var LOG_DEBUG func(msg string, fields ...zapcore.Field)
var LOG_INFO func(msg string, fields ...zapcore.Field)
var LOG_ERROR func(msg string, fields ...zapcore.Field)


func Init_Log() {

	configJson, err := ioutil.ReadFile("logger.json")
	if err != nil {
		panic(err)
	}

	var myConfig zap.Config
	if err := json.Unmarshal(configJson, &myConfig); err != nil {
		panic(err)
	}

	Logger, _ = myConfig.Build()

	LOG_DEBUG = Logger.Debug
	LOG_INFO = Logger.Info
	LOG_ERROR = Logger.Error
}
