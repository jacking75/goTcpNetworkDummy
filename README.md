## 목적
네트워크 라이브러리를 개발 할 때 네트워크 라이브러리의 정확성과 성능을 확인하기 위한 것이다. 

## 테스트 항목  
- 1번 단순한 연결(지정한 개수만큼 서버에 연결한다)
- 2번 연결-끊기 반복
- 3번 고정 길이 크기로 Echo 반복
- 4번 변동 길이 크기로 Echo 반복
- 5번 연결-Echo-끊기 반복
- 6번 연결-Echo-끊기(랜덤) 반복
- 7번 연결-Echo-끊기(서버에서) 반복

## 사용 방법
- 소스 코드를 직접 빌드하면 리눅스, Mac에서 사용가능. 
    
더미 클라이언트의 패킷 헤더는 아래와 같다  
```
int16 size; //헤더와 보디를 포함한 전체 길이
int16 packetId; // 패킷 Id
```  
  
Echo 패킷의 packetId는 101 번이다.  
테스트 항목 7을 위해 더미 클라이언트에서 서버로 접속을 끊어달라고 요청하는 패킷을 보내는데 이때의 packetId는 16번이다(body는 없다).  
    
### 설정 파일 
파일 로그 출력을 위해 logger.json과 더미 클라이언트 동작 설정을 위한 .env  파일이 실행 파일과 같은 디렉토리에 있어야 한다.
  
#### logger.json
아래 예는 debug 레벨에서 로그를 출력한다. 레벨을 변동하고 싶다면 "level" 항목의 값을 바꾼다.  
아래 예는 콘솔 출력만 한다.     
```
{
  "level": "debug",
  "encoding": "json",
  "encoderConfig": {
    "messageKey": "Msg",
    "levelKey": "Level",
    "timeKey": "Time",
    "nameKey": "Name",
    "callerKey": "Caller",
    "stacktraceKey": "St",
    "levelEncoder": "capital",
    "timeEncoder": "iso8601",
    "durationEncoder": "string",
    "callerEncoder": "short"
  },
  "outputPaths": [
    "stdout"
  ],
  "errorOutputPaths": [
    "stderr"
  ]
}
``` 
  
파일 출력을 하고 싶다면 "outputPaths" 항목에 출력을 원하는 파일 이름을 추가한다. 파일 이름은 자유이다.    
```  
"outputPaths": [
    "stdout", "dummylog.log"
  ]
```
  
#### cmd argument
실행 예(Windows)  
```
goTcpNetDummy.exe -c_remoteAddress=127.0.0.1:11021 -c_dummyCount=64 -c_testCase=7 -c_testCountPerDummy=64 -c_testTimeSecondPerDummy=0 -c_sendDataKindCount=8 -c_minSendData=20 -c_maxSendData=50 -c_echoConnectDisconnectRandomPer=50 -c_echoConnectDisconnectServerRandomPer=50
```
    
- c_remoteAddress : 접속할 서버 주소
- c_dummyCount : 더미 개수
- c_testCase : 위의 테스트 항목의 번호에 대응한다
- c_testCountPerDummy : 반복 기준을 횟수로 한다. 
- c_testTimeSecondPerDummy: 반복 기준을 시간(초 단위)으로 한다. 0 보다 크면 시간 기준을 우선한다.
- c_sendDataKindCount : Echo에 사용할 더미 데이터 종류 수
- c_minSendData : Echo에 사용할 데이터의 최소 길이
- c_maxSendData : Echo에 사용할 데이터의 최대 길이
- c_echoConnectDisconnectRandomPer : Echo 하면서 연결을 랜덤하게 끊는 경우 끊을 확률(1~100)
- c_echoConnectDisconnectServerRandomPer : Echo 하면서 서버에서 연결을 랜덤하게 끊는 경우 끊을 확률(1~100)  
  
