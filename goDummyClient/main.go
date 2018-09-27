package main

import (
	"./dummy"
	"./utils"
)

func main() {
	utils.Init_Log()
	utils.Logger.Info("----------- dummy Client Test -----------")

	dummy.Start();
}


