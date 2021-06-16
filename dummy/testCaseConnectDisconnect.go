// 접속만 하고 대기. 몇개까지 접속 가능한지 테스트
package dummy

import (
	"fmt"
	"sync/atomic"

	"go.uber.org/zap"

	"main/utils"
)


func (tester *dummyManager) start_ConnectDisconnect() {
	utils.Logger.Info("start_ConnectDisconnect")
	utils.Logger.Info("",zap.Int("DummyCount", tester.config.dummyCount))

	defer utils.PrintPanicStack()

	for i := range tester.dummyList {
		go tester.DoGoroutine(i)
	}

	go tester.DoGoroutineCheckResult()
}

func (tester *dummyManager) DoGoroutine(dummyIndex int) {
	defer utils.PrintPanicStack()

	config := tester.config

	var remoteAddress string

	if config.isPortByDummy == false {
		remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, config.remotePort)
	} else {
		port := dummyIndex / config.samePortDummyCount
		remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, (config.remotePort + port))
	}

	result := tester.dummyList[dummyIndex].connectAndFailthenSleep(remoteAddress)

	if result {
		atomic.AddInt64(&tester.successCount, 1)
	} else {
		atomic.AddInt64(&tester.failCount, 1)
	}
}




