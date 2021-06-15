package dummy

import (
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"main/utils"
)

// 더미
type dummyObject struct {
	index int			// 고유 인덱스
	conn *net.TCPConn

	repeatConnectCount int64 // 반복 접속 횟수

	echoCount int64 // 에코한 횟수

	recvBuffer []byte // 네트워크 받기 버퍼

	sendPacketQueue *utils.Deque
}

// 더미의 이름 문자열
func (dummy *dummyObject) nameToString() string {
	return fmt.Sprintf("[DummyId: %d]", dummy.index)
}


// 더미 매니져
type dummyManager struct {
	config dummytestConfig
	dummyList []*dummyObject // 테스트에 사용할 더미들

	chDoneEnd chan struct{} // 테스트 완료 통보 채널
	isStop bool // 테스트 중단 여부. 테스트 중이라면 true

	sendDataList [][]byte // 에코에서 사용할 보내기용 데이터 리스트

	// 총 연결 횟수
	connectCount int64
	connectFailCount int64

	// 에코 총 횟수
	echoCount int64

	// 더미 성공, 살패 횟수
	successCount int64
	failCount int64
}



func init_dummyManager(config dummytestConfig) *dummyManager {
	_configValueWriteLogger(config)
	defer utils.PrintPanicStack()

	tester := new(dummyManager)
	tester.config = config

	if checkConfigData(tester) != DUMMY_TEST_ERROR_NONE {
		utils.Logger.Error("Test Config Value Invalide")
		return nil
	}


	rand.Seed(time.Now().UnixNano())

	tester.dummyList = make([]*dummyObject, config.dummyCount)

	for i := range tester.dummyList {
		tester.dummyList[i] = new(dummyObject)
		tester.dummyList[i].index = i

		if config.maxSendData > 0 {
			tester.dummyList[i].recvBuffer = make([]byte, (config.maxSendData+32))
		}

		if config.testCase == TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE ||
			config.testCase == TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE {
			tester.dummyList[i].sendPacketQueue = utils.NewCappedDeque(256)
		}
	}

	tester.chDoneEnd = make(chan struct{}, 1)

	return tester
}

func (tester *dummyManager) DoGoroutineCheckResult() {
	defer utils.PrintPanicStack()

	testStartTime := time.Now()
	testType := tester.config.testCase
	var dummyCount int64 = int64(tester.config.dummyCount)

	for {
		connectCount := atomic.LoadInt64(&tester.connectCount)
		connectFailCount := atomic.LoadInt64(&tester.connectFailCount)
		echoCount := atomic.LoadInt64(&tester.echoCount)
		success := atomic.LoadInt64(&tester.successCount)
		fail := atomic.LoadInt64(&tester.failCount)

		if (success + fail) < dummyCount {
			time.Sleep(1000)
			continue
		}

		utils.Logger.Info("test Completed !!!")

		switch testType {
		case TEST_TYPE_SIMPLE_CONNECT_DISCONNECT:
			utils.Logger.Info("TEST_TYPE_SIMPLE_CONNECT_DISCONNECT", zap.Int64("Connect Count", connectCount), zap.Int64("Connect fail Count", connectFailCount))
		case TEST_TYPE_SIMPLE_REPEAT_CONNECT_DISCONNECT:
			utils.Logger.Info("TEST_TYPE_SIMPLE_REPEAT_CONNECT_DISCONNECT", zap.Int64("Connect Count", connectCount), zap.Int64("Connect fail Count", connectFailCount))
		case TEST_TYPE_ECHO_FIXED_DATA_SIZE:
			utils.Logger.Info("TEST_TYPE_ECHO_FIXED_DATA_SIZE")
		case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE:
			utils.Logger.Info("TEST_TYPE_ECHO_VARIABLE_DATA_SIZE")
		case TEST_TYPE_ECHO_CONNECT_DISCONNECT:
			utils.Logger.Info("TEST_TYPE_ECHO_CONNECT_DISCONNECT")
		case TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM:
			utils.Logger.Info("TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM")
		case TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER:
			utils.Logger.Info("TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER")
		case TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE:
			utils.Logger.Info("TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE")
		case TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE:
			utils.Logger.Info("TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE")
		}

		testTime := time.Now().Sub(testStartTime)
		utils.Logger.Info("test Completed", zap.Int64("Echo Count", echoCount), zap.Int64("Success Count", success), zap.Int64("Fail Count", fail), zap.Duration("Time spent", testTime))

		tester.chDoneEnd <- struct{}{}
		return
	}
}

func (tester *dummyManager) end() {
	utils.Logger.Info("all dummy end")
	defer utils.PrintPanicStack()

	for i := range tester.dummyList {
		if tester.dummyList[i].conn != nil {
			socketClose(tester.dummyList[i])
		}
	}
}


