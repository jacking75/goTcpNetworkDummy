package dummy

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"goTcpNetworkDummy/utils"
)

// 더미 테스트에 사용할 설정 정보
type dummytestConfig struct {
	remoteAddress string // 접속할 서버 주소
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
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Error("Error loading .env file")
	}

	config := dummytestConfig{}
	config.remoteAddress = os.Getenv("remoteAddress")
	config.dummyCount, _ = strconv.Atoi(os.Getenv("dummyCount"))
	config.testCase, _ = strconv.Atoi(os.Getenv("testCase"))
	config.testCountPerDummy, _ = strconv.ParseInt(os.Getenv("testCountPerDummy"), 10, 64)
	config.testTimeSecondPerDummy, _ = strconv.ParseInt(os.Getenv("testTimeSecondPerDummy"), 10, 64)
	config.sendDataKindCount, _ = strconv.Atoi(os.Getenv("sendDataKindCount"))
	config.minSendData, _ = strconv.Atoi(os.Getenv("minSendData"))
	config.maxSendData, _ = strconv.Atoi(os.Getenv("maxSendData"))
	config.echoConnectDisconnectRandomPer, _ = strconv.Atoi(os.Getenv("echoConnectDisconnectRandomPer"))
	config.echoConnectDisconnectServerRandomPer, _ = strconv.Atoi(os.Getenv("echoConnectDisconnectServerRandomPer"))

	return config
}

// 더미 테스트 설정 정보를 출력한다
func _configValueWriteLogger(config dummytestConfig) {
	utils.Logger.Info("init_dummyManager")
	utils.Logger.Info("config", zap.String("Server Address", config.remoteAddress))
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
	case TEST_TYPE_ECHO_FIXED_DATA_SIZE:
		{
			if config.minSendData != config.maxSendData {
				utils.Logger.Error("Echo", zap.Int("MinSize", config.minSendData), zap.Int("MaxSize", config.maxSendData))
				return DUMMY_TEST_ERROR_ECHO_DATA_SIZE
			}
		}
	case TEST_TYPE_ECHO_VARIABLE_DATA_SIZE:
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