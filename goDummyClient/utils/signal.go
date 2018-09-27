package utils

import (
	"os"
	"os/signal"
	"syscall"
)


// handle unix signals
func SignalsHandler(chDoneEnd chan struct{}) bool {
	defer PrintPanicStack()

	Logger.Info("signalsHandler Setting")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	Logger.Info("Waitung.....")

	for {
		select {
			case msg := <-ch:
				{

					switch msg {
					case syscall.SIGTERM: // os 명령어 kill로 종료 시켰음
						Logger.Info("sigterm received: syscall.SIGTERM")
					case syscall.SIGINT: // ctrl + c 로 종료 시켰음
						Logger.Info("sigterm received: syscall.SIGINT")
					}

					Logger.Info("Exit sighandler")
					close(ch)
					return false
				}
			case <-chDoneEnd:
				{
					Logger.Info("Programe Terminate")
					close(ch)
					return true
				}
		}
	}

	return true
}
