package utils

import (
	"go.uber.org/zap"
	//"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
)

var (
	Logger, _ = zap.NewProduction()
)

/*var (
	Logger = logrus.New()
)*/

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
