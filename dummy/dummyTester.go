package dummy

import (
	"time"

	"go.uber.org/zap"

	"main/utils"
)

func Start() {
	config := loadConfig()
	tester := init_dummyManager(config)

	switch config.testCase {
	case TEST_TYPE_SIMPLE_CONNECT_DISCONNECT:
		tester.start_ConnectDisconnect()
	case TEST_TYPE_SIMPLE_REPEAT_CONNECT_DISCONNECT:
		tester.start_RepeatConnectDisconnect()
	case TEST_TYPE_ECHO_FIXED_DATA_SIZE, TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE, TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM:
		fallthrough
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER:
		tester.start_Echo()
	default:
		{
			utils.Logger.Error("Unknown Test Case", zap.Int("case",config.testCase))
			return
		}
	}


	isComplete := utils.SignalsHandler(tester.chDoneEnd)
	tester.isStop = true

	if isComplete == false {
		// 테스트가 중단된 경우라면 테스트 종료를 잠시 기다려 준다.
		time.Sleep(2000)
	}

	tester.end()
}

