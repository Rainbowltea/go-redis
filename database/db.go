package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/resp"
)

type DB struct {
	index int
	data  dict.Dict
}

//所有redis的指令，进行一个实现，
type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}
