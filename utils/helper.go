package utils

import (
	"encoding/binary"
	"errors"
	"runtime"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
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


type RawPacketData struct {
	pos   int
	data  []byte
	order binary.ByteOrder
}

func MakeReader(buffer []byte, isLittleEndian bool) RawPacketData {
	if isLittleEndian {
		return RawPacketData{data: buffer, order: binary.LittleEndian}
	}
	return RawPacketData{data: buffer, order: binary.BigEndian}
}

func MakeWriter(buffer []byte, isLittleEndian bool) RawPacketData {
	if isLittleEndian {
		return RawPacketData{data: buffer, order: binary.LittleEndian}
	}
	return RawPacketData{data: buffer, order: binary.BigEndian}
}

func (p *RawPacketData) Data() []byte {
	return p.data
}

func (p *RawPacketData) Length() int {
	return len(p.data)
}

//=============================================== Readers
func (p *RawPacketData) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

func (p *RawPacketData) ReadS8() (ret int8, err error) {
	_ret, _err := p.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

func (p *RawPacketData) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read uint16 failed")
		return
	}
	buf := p.data[p.pos : p.pos+2]
	ret = p.order.Uint16(buf)
	p.pos += 2
	return
}

func (p *RawPacketData) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *RawPacketData) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = p.order.Uint32(buf)
	p.pos += 4
	return
}

func (p *RawPacketData) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *RawPacketData) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read uint64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = p.order.Uint64(buf)
	p.pos += 8
	return
}

func (p *RawPacketData) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *RawPacketData) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *RawPacketData) ReadBytes(readSize int) (refSlice []byte) {
	refSlice = p.data[p.pos : p.pos+readSize]
	p.pos += readSize
	return
}

func (p *RawPacketData) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string Data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

//================================================ Writers
func (p *RawPacketData) WriteS8(v int8) {
	p.data[p.pos] = (byte)(v)
	p.pos++
}

func (p *RawPacketData) WriteU16(v uint16) {
	p.order.PutUint16(p.data[p.pos:], v)
	p.pos += 2
}

func (p *RawPacketData) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

func (p *RawPacketData) WriteBytes(v []byte) {
	copy(p.data[p.pos:], v)
	p.pos += len(v)
}

func (p *RawPacketData) WriteU32(v uint32) {
	p.order.PutUint32(p.data[p.pos:], v)
	p.pos += 4
}

func (p *RawPacketData) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

func (p *RawPacketData) WriteU64(v uint64) {
	p.order.PutUint64(p.data[p.pos:], v)
	p.pos += 8
}

func (p *RawPacketData) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

func (p *RawPacketData) WriteString(v string) {
	copyLen := copy(p.data[p.pos:], v)
	p.pos += copyLen
}
