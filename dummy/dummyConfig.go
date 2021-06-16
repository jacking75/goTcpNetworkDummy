package dummy

import (
	"flag"
	"go.uber.org/zap"

	"main/utils"
)

// 더미 테스트에 사용할 설정 정보
type dummytestConfig struct {
	remoteIP string // 접속할 서버 IP
	remotePort int // 접속할 서버 Port
	isPortByDummy bool // 더미별로 포트 번호를 별도
	samePortDummyCount int // isPortByDummy를 사용할 때 같은 포트에 접속할 더미 수

	dummyCount int // 더미 개수

	testCase int // 테스트 타입
	testCountPerDummy int64 //테스트 완료 조건 횟수(더미 당)
	testTimeSecondPerDummy int64 // 테스트 완료 조건 시간(초)

	sendDataKindCount int // 에코 테스트 시 보낼 데이터의 종류 수
	minSendData int // 에코 테스트 시 보낼 데이터의 최소 크기. 최소 8 이상
	maxSendData int // 에코 테스트 시 보낼 데이터의 최대 크기. 고정 크기에서는 최대를 사용한다

	echoConnectDisconnectRandomPer int // 에코 연결-끊기 랜덤에서 확률. 100%
	echoConnectDisconnectServerRandomPer int // 에코 연결-끊기(서버가) 랜덤에서 확률. 100%
}

func loadConfig() dummytestConfig {
	config := dummytestConfig{}
	flag.StringVar(&config.remoteIP,"c_remoteIP", "127.0.0.1", "string flag")
	flag.IntVar(&config.remotePort,"c_remotePort", 11021, "int flag")
	flag.BoolVar(&config.isPortByDummy,"c_isPortByDummy", false, "bool flag")
	flag.IntVar(&config.samePortDummyCount,"c_samePortDummyCount", 1, "int flag")

	flag.IntVar(&config.dummyCount,"c_dummyCount", 0, "int flag")
	flag.IntVar(&config.testCase,"c_testCase", 0, "int flag")
	flag.Int64Var(&config.testCountPerDummy,"c_testCountPerDummy", 0, "int flag")
	flag.Int64Var(&config.testTimeSecondPerDummy,"c_testTimeSecondPerDummy", 0, "int flag")
	flag.IntVar(&config.sendDataKindCount,"c_sendDataKindCount", 0, "int flag")
	flag.IntVar(&config.minSendData,"c_minSendData", 0, "int flag")
	flag.IntVar(&config.maxSendData,"c_maxSendData", 0, "int flag")
	flag.IntVar(&config.echoConnectDisconnectRandomPer,"c_echoConnectDisconnectRandomPer", 0, "int flag")
	flag.IntVar(&config.echoConnectDisconnectServerRandomPer,"c_echoConnectDisconnectServerRandomPer", 0, "int flag")

	flag.Parse()

	return config
}

// 더미 테스트 설정 정보를 출력한다
func _configValueWriteLogger(config dummytestConfig) {
	utils.Logger.Info("init_dummyManager")
	utils.Logger.Info("config", zap.String("Server IP", config.remoteIP))
	utils.Logger.Info("config", zap.Int("Server Port: ", config.remotePort))
	utils.Logger.Info("config", zap.Bool("isPortByDummy: ", config.isPortByDummy))
	utils.Logger.Info("config", zap.Int("DummyCount: ", config.dummyCount))
	utils.Logger.Info("config", zap.Int("Test Case: ", config.testCase))
	utils.Logger.Info("config", zap.Int64("testCountPerDummy: ", config.testCountPerDummy))
	utils.Logger.Info("config", zap.Int64("testTimeSecondPerDummy: ", config.testTimeSecondPerDummy))
	utils.Logger.Info("config", zap.Int("sendDataKindCount: ", config.sendDataKindCount))
	utils.Logger.Info("config", zap.Int("minSendData: ", config.minSendData))
	utils.Logger.Info("config", zap.Int("maxSendData: ", config.maxSendData))
	utils.Logger.Info("config", zap.Int("echoConnectDisconnectRandomPer: ", config.echoConnectDisconnectRandomPer))
	utils.Logger.Info("config", zap.Int("echoConnectDisconnectServerRandomPer: ", config.echoConnectDisconnectServerRandomPer))
}

// 더미 테스트 설정 값이 올바른지 조사한다
func checkConfigData(tester *dummyManager) int {
	config := tester.config

	if config.minSendData < 8 {
		utils.Logger.Error("Echo", zap.Int("MinSize", config.minSendData))
		return DUMMY_TEST_ERROR_ECHO_DATA_MIN_SIZE
	}

	switch config.testCase {
	case TEST_TYPE_ECHO_FIXED_DATA_SIZE, TEST_TYPE_ECHO_EX_FIXED_DATA_SIZE:
		{
			if config.minSendData != config.maxSendData {
				utils.Logger.Error("Echo", zap.Int("MinSize", config.minSendData), zap.Int("MaxSize", config.maxSendData))
				return DUMMY_TEST_ERROR_ECHO_DATA_SIZE
			}
		}
	case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE, TEST_TYPE_ECHO_EX_VARIABLE_DATA_SIZE:
		{
			if config.minSendData == config.maxSendData {
				utils.Logger.Error("Echo", zap.Int("MinSize", config.minSendData), zap.Int("MaxSize", config.maxSendData))
				return DUMMY_TEST_ERROR_ECHO_DATA_SIZE
			}
		}
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM:
		{
			if config.echoConnectDisconnectRandomPer <= 0 ||
				config.echoConnectDisconnectRandomPer > 100 {
				utils.Logger.Error("Invalide echoConnectDisconnectRandomPer", zap.Int("echoConnectDisconnectRandomPer", config.echoConnectDisconnectRandomPer))
				return DUMMY_TEST_ERROR_ECHO_CONNECT_DISCONNECT_RANDOMPER
			}
		}
	case TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER:
		{
			if config.echoConnectDisconnectServerRandomPer <= 0 ||
				config.echoConnectDisconnectServerRandomPer > 100 {
				utils.Logger.Error("Invalide echoConnectDisconnectServerRandomPer", zap.Int("echoConnectDisconnectServerRandomPer",config.echoConnectDisconnectServerRandomPer))
				return DUMMY_TEST_ERROR_ECHO_CONNECT_DISCONNECT_SERVER_RANDOMPER
			}
		}
	}

	return DUMMY_TEST_ERROR_NONE
}