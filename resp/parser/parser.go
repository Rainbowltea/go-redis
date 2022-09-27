package parser

import (
	"bufio"
	"errors"
	"go-redis/interface/resp"
	"io"
)

// Payload stores redis.Reply or error
type Payload struct {
	Data resp.Reply
	Err  error
}
type readState struct {
	readingMultiLine  bool     //记录需要解析单行指令还是多行指令
	expectedArgsCount int      //应该解析多少个参数
	msgType           byte     //用户消息
	args              [][]byte //已经解析的参数  SET ..   KEY  ..  Val
	bulkLen           int64    //长度
}

//解析是否完成
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

//异步解析，业务处理和协议解析并发运行
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}
func parse0(reader io.Reader, ch chan<- *Payload) {

}

//读进数据
//1.按\r\n切分
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 { // read normal line
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else { // read bulk line (binary safe)
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 ||
			msg[len(msg)-2] != '\r' ||
			msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}
