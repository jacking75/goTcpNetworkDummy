// 접속 후 에코 - 데이터 크기 고정
package dummy

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/kyokomi/lottery"
	"go.uber.org/zap"

	"main/utils"
)

func (tester *dummyManager) start_Echo() {
	utils.Logger.Info("start_Echo")
	defer utils.PrintPanicStack()

	config := tester.config

	switch config.testCase {
	case TEST_TYPE_ECHO_FIXED_DATA_SIZE, TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE, TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM:
		tester.sendDataList = makeRandomPacketList(config.sendDataKindCount, config.minSendData, config.maxSendData)
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER:
		tester.sendDataList = makePackets_Normal_ReqDisConn_Both(config.minSendData, config.maxSendData)
	}

	utils.Logger.Info("start_Echo - start goroutine")
	for i := range tester.dummyList {
		var remoteAddress string

		if config.isPortByDummy == false {
			remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, config.remotePort)
		} else {
			port := i / config.samePortDummyCount
			remoteAddress = fmt.Sprintf("%s:%d", config.remoteIP, (config.remotePort + port))
		}
		go tester._DoGoroutine_Echo(i, remoteAddress, config.testCase)
	}

	go tester.DoGoroutineCheckResult()
}

func (tester *dummyManager) _DoGoroutine_Echo(dummyIndex int, remoteAddress string, testCase int) {
	//utils.Logger.Debug("_DoGoroutine_Echo")
	defer utils.PrintPanicStack()

	result := tester._Echo(dummyIndex, remoteAddress, testCompleteCheckDoCountOrTime)

	if result {
		atomic.AddInt64(&tester.successCount, 1)
	} else {
		atomic.AddInt64(&tester.failCount, 1)
	}

	// 에코인데 받지 않고 막 보내는 타입은 작업을 양보한다. 그렇지 않고 막보내면 패킷이 뭉쳐서 서버로 보내지게될 확률이 높다
	if testCase == TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE || testCase == TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE {
		runtime.Gosched()
	}
}


func (tester *dummyManager) _Echo(dummyIndex int, remoteAddress string,
										in_completeCheck func(bool,int64,int64) bool) bool {
	//utils.Logger.Debug("_Echo")
	defer utils.PrintPanicStack()

	dummy := tester.dummyList[dummyIndex]
	config := tester.config
	isCountCheck := config.testCountPerDummy > 0
	limitCount := tester.config.testCountPerDummy

	if isCountCheck == false {
		limitCount = time.Now().Unix() + tester.config.testTimeSecondPerDummy
	}


	lot1 := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano ())))

	var serverDisconnFlag int16 // 랜덤하게 서버에서 짜르는 것.
	serverDisconnRandom := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano ())))

	for {
		sendData := tester._selectPacket(config, serverDisconnRandom, serverDisconnFlag)
		var result int

		if dummy.sendPacketQueue == nil {
			result = dummy.connectAndEcho(remoteAddress, sendData)
		} else {
			result = dummy.connectAndEchoWithoutReceive(remoteAddress, sendData)
		}

		if _checkErrorEnableContinue(config.testCase, result) == false {
			utils.Logger.Info("_Echo", zap.Int("Error", result))
			return false
		}


		atomic.AddInt64(&tester.echoCount, 1)
		dummy.echoCount++

		if in_completeCheck(isCountCheck, limitCount, dummy.echoCount) {
			return true
		}

		if tester.isStop == true {
			return true
		}

		_testCaseCheckThenDisconnect(config, dummy, lot1)
	}

	return true
}

// 에러 코드를 조사하고, 테스트를 진행할 수 없다면 false 반환
func _checkErrorEnableContinue(testCase int, errorCode int) bool {
	if errorCode == NET_ERROR_NONE {
		return true
	}

	if errorCode == NET_ERROR_ERROR_DISCONNECTED && testCase == TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER {
		return true
	}

	return false
}

// 에코로 어떤 데이터를 보낼지 선택한다
func (tester *dummyManager) _selectPacket(config dummytestConfig, lot lottery.Lottery, serverDisconnFlag int16) [] byte {
	testCase := config.testCase
	sendDataKindCount := config.sendDataKindCount
	var sendData []byte

	switch testCase {
	case TEST_TYPE_ECHO_FIXED_DATA_SIZE, TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE, TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM:
		{
			sendIndex := global_randNumber(sendDataKindCount)
			if sendIndex < 0 {
				sendIndex = 0
			} else if sendIndex >= sendDataKindCount {
				sendIndex = sendDataKindCount
			}
			sendData = tester.sendDataList[sendIndex]
		}
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER:
		{
			if lot.Lot ( config.echoConnectDisconnectServerRandomPer ) {
				if serverDisconnFlag == 0 {
					sendData = tester.sendDataList[1]
					serverDisconnFlag = 1
				} else {
					sendData = tester.sendDataList[2]
					serverDisconnFlag = 0
				}

			} else {
				sendData = tester.sendDataList[0]
			}
		}
	default:
		{
			utils.Logger.Error("Unknown Test Case", zap.Int("Case", testCase))
		}
	}

	return sendData
}

// 테스트 중 연결을 끊어야 하는 경우라면 연결을 끊는다
func _testCaseCheckThenDisconnect(config dummytestConfig, dummy *dummyObject, lot lottery.Lottery) {
	isClose := false

	if config.testCase == TEST_TYPE_ECHO_CONNECT_DISCONNECT {
		isClose = true
	} else if config.testCase == TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM {
		if lot.Lot ( config.echoConnectDisconnectRandomPer ) {
			isClose = true
		}
	}

	if isClose {
		socketClose(dummy)
		//utils.Logger.Debug("_testCaseCheckThenDisconnect")
	}
}