package main

import (
	"github.com/duhaifeng/loglet"
	"github.com/duhaifeng/loglet/call"
	"time"
)

func main() {
	call.M1()

	logger := loglet.NewLogger()
	logger.Debug("111111111111111111111111111")
	time.Sleep(time.Second)
}
