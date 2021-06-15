package main

import (
	"main/dummy"
	. "main/utils"
)

func main() {
	Init_Log()

	LOG_INFO("----------- dummy Client Test -----------")

	dummy.Start()
}


