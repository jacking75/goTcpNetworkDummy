package dummy

import (
	"time"
	"net"

	"go.uber.org/zap"

	"goTcpNetworkDummy/utils"
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

// 에코 - 접속 후 데이터를 보내고 받는다
func (dummy *dummyObject) connectAndEcho(remoteAddress string, sendData []byte) int {
	utils.Logger.Debug("connectAndEcho Start", zap.String("Dummy", dummy.nameToString()))

	if dummy.conn == nil {
		utils.Logger.Debug("ConnectAndEcho. Connect", zap.String("Dummy",dummy.nameToString()))
		for {
			if dummy.connectAndFailthenSleep(remoteAddress) == false {
				//utils.Logger.Debug("ConnectAndEcho. fail")
				continue
			}
			break
		}
	}
	//utils.Logger.Debug("ConnectAndEcho. Complete")

	sendSize := len(sendData)
	writeBytes, err1 := dummy.conn.Write(sendData)
	//utils.Logger.Debug("ConnectAndEcho. Write")
	if err1 != nil {
		utils.Logger.Error("ConnectAndEcho", zap.String("Dummy", dummy.nameToString()), zap.Error(err1))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND
	}

	if writeBytes != sendSize {
		utils.Logger.Error("Tcp Write Diff Size. Expact-Result", zap.String("Dummy",dummy.nameToString()), zap.Int("sendSize",sendSize), zap.Int("writeBytes", writeBytes))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND_DIFF_SIZE
	}

	//utils.Logger.Debug("ConnectAndEcho. read start")
	recvBytes, err2 := dummy.conn.Read(dummy.recvBuffer)
	//utils.Logger.Debug("ConnectAndEcho. read end")
	if recvBytes == 0 {
		socketClose(dummy)
		return NET_ERROR_ERROR_DISCONNECTED
	}

	if err2 != nil {
		utils.Logger.Error("Tcp Read error", zap.String("Dummy", dummy.nameToString()), zap.Error(err2))
		socketClose(dummy)
		return NET_ERROR_ERROR_RECV
	}


	//보낸 데이터가 그대로 왔는지 확인
	if sendSize != recvBytes {
		utils.Logger.Error("Tcp Read Diff Size. Expact-Result", zap.String("Dummy",dummy.nameToString()), zap.Int("sendSize",sendSize), zap.Int("recvBytes", recvBytes))
		return NET_ERROR_ERROR_SEND_RECV_DIFF_SIZE
	}

	if sendData[8] != dummy.recvBuffer[8] || sendData[8] != dummy.recvBuffer[8] {
		return NET_ERROR_ERROR_SEND_RECV_DIFF_DATA
	}

	//TODO 보낸 데이터와 받는 데이터가 같은지 검증하기
	//utils.Logger.Debug("connectAndEcho. send-receive data Size:", sendSize)
	return NET_ERROR_NONE
}



// 리모트 컴퓨터에 접속한다
func (dummy *dummyObject) _remoteConnect(address string) bool {
	//utils.Logger.Debug("_remoteConnect :", address)
	defer utils.PrintPanicStack()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		utils.Logger.Error("fail ResolveTCPAddr address", zap.String("Dummy", dummy.nameToString()), zap.Error(err))
		return false
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		utils.Logger.Error("fail DialTCP", zap.String("Dummy", dummy.nameToString()), zap.Error(err))
		return false
	}

	dummy.conn = conn
	dummy.conn.SetLinger(0)
	//utils.Logger.Info("success: _remoteConnect")
	return true
}