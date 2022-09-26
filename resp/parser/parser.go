package parser

import (
	"go-redis/interface/resp"
)

// Payload stores redis.Reply or error
type Payload struct {
	Data resp.Reply
	Err  error
}
type readState struct {
	readingMultiLine  bool     //记录需要解析单行指令还是多行指令
	expectedArgsCount int      //记录需要解析有多少个参数
	msgType           byte     //用户消息
	args              [][]byte //用户传入的数据本身
	bulkLen           int64    //长度
}
