package utils

import (
	"encoding/binary"
	"log"
)

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Int64ToByte(num int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}
