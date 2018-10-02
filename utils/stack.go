package utils

import (
	"runtime"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"fmt"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		Logger.Error(fmt.Sprintf("%v", x))

		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			Logger.Error("frame", zap.Int("N", i), zap.String("func",runtime.FuncForPC(funcName).Name()), zap.String("file", file), zap.Int("line", line))
			//Logger.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			msg := fmt.Sprintf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
			Logger.Error(msg)
		}
	}
}
