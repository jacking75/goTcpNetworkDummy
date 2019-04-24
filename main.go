package main

import (
	"goTcpNetworkDummy/dummy"
	. "goTcpNetworkDummy/utils"
)

func main() {
	Init_Log()

	LOG_INFO("----------- dummy Client Test -----------")

	dummy.Start()
}


