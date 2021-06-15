package dummy

import (
	"go.uber.org/zap"
	"net"

	. "main/utils"
)

// 에코 - 접속 후 데이터를 보내고 받는다. 중간에 접속이 끊어지지 않는다.
func (dummy *dummyObject) connectAndEchoWithoutReceive(remoteAddress string, sendData []byte) int {
	LOG_DEBUG("connectAndEchoWithoutReceive Start", zap.String("Dummy", dummy.nameToString()))

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
		go _echoReceive_goroutine(dummy.nameToString(), dummy.conn, dummy.recvBuffer, dummy.sendPacketQueue)
	}

	if _receiveErrorCode != NET_ERROR_NONE {
		socketClose(dummy)
		return _receiveErrorCode
	}


	dummy.sendPacketQueue.Append(sendData)

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

	return NET_ERROR_NONE
}

func _echoReceive_goroutine(dummyName string, conn *net.TCPConn, recvBuffer []byte, sendPacketQueue *Deque) {
	LOG_DEBUG("_echoReceive_goroutine. start")
	recvStartPos := 0

	for {
		readAbleBytes, err := conn.Read(recvBuffer[recvStartPos:])
		//LOG_DEBUG("_echoReceive_goroutine. read end")

		_receiveErrorCode = _checkReadError(dummyName, readAbleBytes, err)
		if _receiveErrorCode != NET_ERROR_NONE {
			conn.Close()
			return
		}


		readAbleBytes += recvStartPos
		recvStartPos, _receiveErrorCode = _checkReceiveData(readAbleBytes, sendPacketQueue, recvBuffer)

		if _receiveErrorCode != NET_ERROR_NONE {
			conn.Close()
			return
		}
	}
}


var _receiveErrorCode int = NET_ERROR_NONE

func _checkReadError(dummyName string, readAbleBytes int, netErr error) int {
	if readAbleBytes == 0 {
		return NET_ERROR_ERROR_DISCONNECTED
	}

	if netErr != nil {
		LOG_ERROR("Tcp Read error", zap.String("Dummy", dummyName), zap.Error(netErr))
		return NET_ERROR_ERROR_RECV
	}

	return NET_ERROR_NONE
}

func _checkReceiveData(readAbleBytes int, sendPacketQueue *Deque, recvBuffer []byte) (int, int) {
	readBufferPos := 0
	recvStartPos := 0
	errorCode := NET_ERROR_NONE

	for {
		if readAbleBytes <= 0 {
			break
		}

		sendPacket := sendPacketQueue.Shift().([]byte)
		sendPacketSize := len(sendPacket)

		if sendPacketSize > readAbleBytes {
			sendPacketQueue.Prepend(sendPacket)
			break
		}

		if sendPacket[0] != recvBuffer[readBufferPos] ||
			sendPacket[sendPacketSize-1] != recvBuffer[readBufferPos+sendPacketSize-1] {
			return recvStartPos, NET_ERROR_ERROR_DISCONNECTED
		}

		readAbleBytes -= sendPacketSize
		readBufferPos += sendPacketSize
	}

	if readAbleBytes > 0 {
		copy(recvBuffer, recvBuffer[readBufferPos:(readBufferPos+ readAbleBytes)])
	}

	recvStartPos = readAbleBytes
	return recvStartPos, errorCode
}

