package tool

import (
	"time"
)

func Rand(s int64) int64 {
	if s == 0 {
		return 0
	}
	return time.Now().Unix() % s
}
