// 반복해서 접속과 연결종료
package dummy

import (
	"fmt"
	"sync/atomic"
	"time"

	"main/utils"
)


func (tester *dummyManager) start_RepeatConnectDisconnect() {
	utils.Logger.Info("start_RepeatConnectDisconnect")
	defer utils.PrintPanicStack()

	config := tester.config
	isCountCheck := config.testCountPerDummy > 0

	for i := range tester.dummyList {
		var remoteAddress string

		if config.isPortByDummy == false {
			remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, config.remotePort)
		} else {
			port := i / config.samePortDummyCount
			remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, (config.remotePort + port))
		}

		go tester._DoGoroutine_RepeatConnectDisconnect(i, remoteAddress, isCountCheck)
	}

	go tester.DoGoroutineCheckResult()
}

func (tester *dummyManager) _DoGoroutine_RepeatConnectDisconnect(dummyIndex int, remoteAddress string, isCountCheck bool) {
	defer utils.PrintPanicStack()

	result := tester._repeatConnectDisconnec(dummyIndex, remoteAddress, testCompleteCheckDoCountOrTime)

	if result {
		atomic.AddInt64(&tester.successCount, 1)
	}
}

func (tester *dummyManager) _repeatConnectDisconnec(dummyIndex int, remoteAddress string,
													completeCheck func(bool,int64,int64) bool) bool {
	//utils.Logger.Info("_repeatConnectDisconnec")

	dummy := tester.dummyList[dummyIndex]
	isCountCheck := tester.config.testCountPerDummy > 0
	var limitCount int64

	if isCountCheck {
		limitCount = tester.config.testCountPerDummy
	} else {
		limitCount = time.Now().Unix() + tester.config.testTimeSecondPerDummy
	}


	for {
		if ret := dummy._connectAndDisconnec(remoteAddress); ret == false {
			atomic.AddInt64(&tester.connectFailCount, 1)
			continue
		}


		atomic.AddInt64(&tester.connectCount, 1)

		if completeCheck(isCountCheck, limitCount, dummy.repeatConnectCount) {
			return true
		}

		if tester.isStop == true {
			return true
		}
	}

	return true
}

func (dummy *dummyObject) _connectAndDisconnec(remoteAddress string) bool {
	//utils.Logger.Info("_connectAndDisconnec")

	if dummy.connectAndFailthenSleep(remoteAddress) == false {
		return false
	}

	socketClose(dummy)
	dummy.repeatConnectCount++

	return true
}