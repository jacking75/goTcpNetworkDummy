package main

import (
	"goTcpNetworkDummy/dummy"
	"goTcpNetworkDummy/utils"
)

func main() {
	utils.Init_Log()
	utils.Logger.Info("----------- dummy Client Test -----------")

	dummy.Start();
}


