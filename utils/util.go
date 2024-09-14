package utils

import (
	"encoding/binary"
	"log"
	"os"
)

// Int64ToByte int64转换为byte数组
func Int64ToByte(num int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num)) //将num转换为uint64后存入buf
	return buf
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
