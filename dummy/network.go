package dummy

import (
	"net"
	"time"

	"go.uber.org/zap"

	. "main/utils"
)

// 네트워크 관련 에커 코드 선언
const (
	NET_ERROR_NONE = 0
	NET_ERROR_ERROR_SEND = 1
	NET_ERROR_ERROR_RECV = 2
	NET_ERROR_ERROR_DISCONNECTED = 3
	NET_ERROR_ERROR_SEND_DIFF_SIZE = 4
	NET_ERROR_ERROR_SEND_RECV_DIFF_SIZE = 5
	NET_ERROR_ERROR_SEND_RECV_DIFF_DATA = 6
)

// 연결을 시도하고 실패하면 대기 후 다시 시도한다.
// 더미 클라이언트 테스트에서 같은 타이밍에 너무 많은 접속 요구가 있으면 서버가 받아들이지 못할 수 있음.
func (dummy *dummyObject) connectAndFailthenSleep(remoteAddress string) bool {
	result := dummy._remoteConnect(remoteAddress)

	if result == false {
		// 아마 서버에서 listen 처리 중이므로 잠깐 대기 후 시도한다
		millisecond := (int64)(global_randNumber(100))
		time.Sleep(time.Duration(millisecond))
		return false
	}

	return true
}

// 소켓 접속을 끊는다
func socketClose(dummy *dummyObject) {
	dummy.conn.Close()
	dummy.conn = nil
}

// 리모트 컴퓨터에 접속한다
func (dummy *dummyObject) _remoteConnect(address string) bool {
	//LOG_DEBUG("_remoteConnect :", address)
	defer PrintPanicStack()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		LOG_ERROR("fail ResolveTCPAddr address", zap.String("Dummy", dummy.nameToString()), zap.Error(err))
		return false
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		LOG_ERROR("fail DialTCP", zap.String("Dummy", dummy.nameToString()), zap.Error(err))
		return false
	}

	dummy.conn = conn
	dummy.conn.SetLinger(0)
	//LOG_INFO("success: _remoteConnect")
	return true
}