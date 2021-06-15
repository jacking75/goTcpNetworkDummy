package dummy

import (
	"go.uber.org/zap"

	. "main/utils"
)

// 에코 - 접속 후 데이터를 보내고 받는다
func (dummy *dummyObject) connectAndEcho(remoteAddress string, sendData []byte) int {
	LOG_DEBUG("connectAndEcho Start", zap.String("Dummy", dummy.nameToString()))

	if dummy.conn == nil {
		LOG_DEBUG("ConnectAndEcho. Connect", zap.String("Dummy",dummy.nameToString()))
		for {
			if dummy.connectAndFailthenSleep(remoteAddress) == false {
				//LOG_DEBUG("ConnectAndEcho. fail")
				continue
			}
			break
		}
	}

	/*packetId := binary.LittleEndian.Uint16(sendData[2:4])
	if packetId == ECHO_REQ_DISCONNECT_PACKET_ID {
		LOG_INFO("send ECHO_REQ_DISCONNECT_PACKET_ID", zap.String("Dummy", dummy.nameToString()))
	}*/

	sendSize := len(sendData)
	writeBytes, err1 := dummy.conn.Write(sendData)
	//LOG_DEBUG("ConnectAndEcho. Write")
	if err1 != nil {
		LOG_ERROR("ConnectAndEcho", zap.String("Dummy", dummy.nameToString()), zap.Error(err1))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND
	}

	if writeBytes != sendSize {
		LOG_ERROR("Tcp Write Diff Size. Expact-Result", zap.String("Dummy",dummy.nameToString()), zap.Int("sendSize",sendSize), zap.Int("writeBytes", writeBytes))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND_DIFF_SIZE
	}

	//LOG_DEBUG("ConnectAndEcho. read start")
	recvBytes, err2 := dummy.conn.Read(dummy.recvBuffer)
	//LOG_DEBUG("ConnectAndEcho. read end")
	if recvBytes == 0 {
		//LOG_DEBUG("Closed. ConnectAndEcho", zap.String("Dummy", dummy.nameToString()))
		socketClose(dummy)
		return NET_ERROR_ERROR_DISCONNECTED
	}

	if err2 != nil {
		LOG_ERROR("Tcp Read error", zap.String("Dummy", dummy.nameToString()), zap.Error(err2))
		socketClose(dummy)
		return NET_ERROR_ERROR_RECV
	}


	//보낸 데이터가 그대로 왔는지 확인
	if sendSize != recvBytes {
		LOG_ERROR("Tcp Read Diff Size. Expact-Result", zap.String("Dummy",dummy.nameToString()), zap.Int("sendSize",sendSize), zap.Int("recvBytes", recvBytes))
		return NET_ERROR_ERROR_SEND_RECV_DIFF_SIZE
	}

	if sendData[0] != dummy.recvBuffer[0] || sendData[sendSize-1] != dummy.recvBuffer[sendSize-1] {
		return NET_ERROR_ERROR_SEND_RECV_DIFF_DATA
	}

	//TODO 보낸 데이터와 받는 데이터가 같은지 검증하기
	LOG_DEBUG("connectAndEcho. send-receive: ", zap.Int("sendSize", sendSize))
	return NET_ERROR_NONE
}

