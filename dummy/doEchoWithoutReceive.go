package dummy

import (
	"dancevilleparkserver/autoTest/dummy"
	"go.uber.org/zap"
	"net"

	. "goTcpNetworkDummy/utils"
)

// 에코 - 접속 후 데이터를 보내고 받는다. 중간에 접속이 끊어지지 않는다.
func (dummy *dummyObject) connectAndEchoWithoutReceive(remoteAddress string, sendData []byte, sendPacketQueue *Deque) int {
	LOG_DEBUG("connectAndEchoWithoutReceive Start", zap.String("Dummy", dummy.nameToString()))

	var errorCode = NET_ERROR_NONE

	if dummy.conn == nil {
		LOG_DEBUG("connectAndEchoWithoutReceive. Connect", zap.String("Dummy",dummy.nameToString()))
		for {
			if dummy.connectAndFailthenSleep(remoteAddress) == false {
				//LOG_DEBUG("ConnectAndEcho. fail")
				continue
			}
			break
		}

		//golang은 Nagle 알고리즘이 기본은 off
		go _echoReceive_goroutine(dummy.nameToString(), dummy.conn, dummy.recvBuffer, sendPacketQueue, &errorCode)

		LOG_DEBUG("connectAndEchoWithoutReceive. Init", zap.String("Dummy",dummy.nameToString()))
	}
	//utils.Logger.Debug("connectAndEchoWithoutReceive. Connect, Init")


	sendSize := len(sendData)
	writeBytes, err1 := dummy.conn.Write(sendData)
	//LOG_DEBUG("connectAndEchoWithoutReceive. Write")
	if err1 != nil {
		LOG_ERROR("connectAndEchoWithoutReceive", zap.String("Dummy", dummy.nameToString()), zap.Error(err1))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND
	}

	if writeBytes != sendSize {
		LOG_ERROR("Tcp Write Diff Size. Expact-Result", zap.String("Dummy",dummy.nameToString()), zap.Int("sendSize",sendSize), zap.Int("writeBytes", writeBytes))
		socketClose(dummy)
		return NET_ERROR_ERROR_SEND_DIFF_SIZE
	}

	sendPacketQueue.Append(sendData)


	if errorCode != NET_ERROR_NONE {
		return errorCode
	}

	return NET_ERROR_NONE
}

func _echoReceive_goroutine(dummyName string, conn *net.TCPConn, recvBuffer []byte, sendPacketQueue *Deque, errorCode *int) {
	//LOG_DEBUG("_echoReceive_goroutine. start")
	recvPos := 0

	for {
		recvBytes, err2 := conn.Read(recvBuffer[recvPos:])
		//LOG_DEBUG("_echoReceive_goroutine. read end")
		if recvBytes == 0 {
			*errorCode = NET_ERROR_ERROR_DISCONNECTED
			conn.Close()
			return
		}

		if err2 != nil {
			LOG_ERROR("Tcp Read error", zap.String("Dummy", dummyName), zap.Error(err2))
			*errorCode = NET_ERROR_ERROR_RECV
			conn.Close()
			return
		}

		recvBytes += recvPos
		readBufferPos := 0

		// 패킷 분해
		for {
			sendPacket := sendPacketQueue.Pop()
			sendPacketSize := len(sendPacket.([]byte))

			if sendPacketSize == recvBytes {
				recvPos = 0
				break
			} else if sendPacketSize < recvBytes {
				recvBytes -= sendPacketSize
				readBufferPos += recvBytes
			} else {
				sendPacketQueue.Append(sendPacket)
				recvPos = recvBytes
				break
			}
		}


		if sendData[8] != dummy.recvBuffer[8] || sendData[8] != dummy.recvBuffer[8] {
			return//NET_ERROR_ERROR_SEND_RECV_DIFF_DATA
		}

		//TODO 보낸 데이터와 받는 데이터가 같은지 검증하기
		//LOG_DEBUG("connectAndEcho. send-receive data Size:", sendSize)
	}
}


