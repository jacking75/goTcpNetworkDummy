package utils

import (
	"io/ioutil"
	"encoding/json"

	"go.uber.org/zap"
)

var (
	Logger, _ = zap.NewProduction()
)


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
}
