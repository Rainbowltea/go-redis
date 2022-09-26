package resp

type Connection interface { //使用接口时考虑到和AOF相关有不同的实现
	Write([]byte) error //回复客户端
	// used for multi database
	GetDBIndex() int
	SelectDB(int)
}
