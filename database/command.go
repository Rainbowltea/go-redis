package database

import "strings"

//记录系统中的所有指令与command的关系
var cmdTable = make(map[string]*command)

type command struct {
	exector ExecFunc
	artiy   int //参数的数量
}

//
func RegisterCommand(name string, exector ExecFunc, artiy int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: exector,
		artiy:   artiy,
	}
}
