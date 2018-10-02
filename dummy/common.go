package dummy

import (
	"go.uber.org/zap"
	"math/rand"
	"time"
	"encoding/binary"
)


const (
	TEST_TYPE_SIMPLE_CONNECT_DISCONNECT = 1 // 단순 연결
	TEST_TYPE_SIMPLE_REPEAT_CONNECT_DISCONNECT = 2 // 반복 연결-끊기
	TEST_TYPE_ECHO_FIXED_DATA_SIZE = 3 // 에코 고정 길이 데이터
	TEST_TYPE_ECHO_VARIABLE_DATA_SIZE = 4// 임의 크기 데이터 에코
	TEST_TYPE_ECHO_CONNECT_DISCONNECT = 5// 연결-에코-끊기
	TEST_TYPE_ECHO_CONNECT_DISCONNECT_RANDOM = 6// 연결-에코-끊기(랜덤)
	TEST_TYPE_ECHO_CONNECT_DISCONNECT_FROM_SERVER = 7// 연결-에코-서버에서 끊기
)

const (
	DUMMY_TEST_ERROR_NONE = 0
	DUMMY_TEST_ERROR_ECHO_DATA_SIZE = 1
	DUMMY_TEST_ERROR_ECHO_DATA_MIN_SIZE = 2
	DUMMY_TEST_ERROR_ECHO_CONNECT_DISCONNECT_RANDOMPER = 3
	DUMMY_TEST_ERROR_ECHO_CONNECT_DISCONNECT_SERVER_RANDOMPER = 4
)


func global_randNumber(maxNumber int) int {
	return rand.Intn(maxNumber)
}

// 테스트 완료를 지정횟수 or 지정 시간으로 조사한다.
func testCompleteCheckDoCountOrTime(isCount bool, limitValue int64, doCount int64) bool {
	if isCount {
		if doCount >= limitValue {
			return true
		}

	} else {
		second := time.Now().Unix()
		if limitValue <= second {
			return true
		}
	}

	return false
}

// 랜덤한 패킷을 지정 크기와, 갯수만큰 만든다. minSize와 maxSize가 같은면 패킷의 크기는 고정 크기
func makeRandomPacketList(row int, minSize int, maxSize int) [][]byte {
	binaryDatas := make([][]byte, row)

	for i := 0; i < row; i++ {
		binaryDatas[i] = _randomBytes(minSize, maxSize)
	}

	utils.Logger.Info("makeRandomPacketList", zap.Int("Row", row), zap.Int("MinSize", minSize), zap.Int("MaxSize", maxSize))
	return binaryDatas
}

// 보통 패킷, 연결 종료 요청 패킷, 보통+연결종료요청 2개 합쳐진 패킷 을 만든다.
func makePackets_Normal_ReqDisConn_Both(size int) [][]byte {
	binaryDatas := make([][]byte, 3)

	binaryDatas[0] = _randomBytes(size, size)
	binaryDatas[1] = _requestDisconnecPacketBytes(16, 4)

	binaryDatas[2] = make([]byte, (len(binaryDatas[0])+len(binaryDatas[1])))
	binaryDatas[2] = append(binaryDatas[2], binaryDatas[0]...)
	binaryDatas[2] = append(binaryDatas[2], binaryDatas[1]...)

	utils.Logger.Info("makePackets_Normal_ReqDisConn_Both", zap.Int("Row", 3), zap.Int("Size", size))
	return binaryDatas
}



func _global_randRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func _randomBytes(minSize int, maxSize int) []byte {
	var size int

	isRandomSize := false
	if minSize != maxSize {
		isRandomSize = true
	}

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if isRandomSize {
		size = _global_randRange(minSize, maxSize)
	} else {
		size = maxSize
	}

	b := make([]byte, size)
	binary.LittleEndian.PutUint16(b[0:2], (uint16)(size))
	binary.LittleEndian.PutUint16(b[2:4], 101)

	for n:=4; n < size; n++{
		b[n] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}

func _requestDisconnecPacketBytes(packetId int16, totalPacketSize int16) []byte {
	b := make([]byte, totalPacketSize)
	binary.LittleEndian.PutUint16(b[0:2], (uint16)(totalPacketSize))
	binary.LittleEndian.PutUint16(b[2:4], (uint16)(packetId))


	bodySize := (int)(totalPacketSize - 4)
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for n:=4; n < bodySize; n++{
		b[n] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}