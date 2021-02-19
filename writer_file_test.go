package loglet

import (
	"strconv"
	"testing"
	"time"
)

const line = "ABCDEFGHIJKLMNOPQRSTUVWXYZ---日志--- "

func TestFileWriter(t *testing.T) {
	logger := new(FileWriter)
	logger.SetFileBaseName("/var/loglet/log_test/test.loglet")
	logger.WriteLog(&LogMsg{msgLevel: DEBUG, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: INFO, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: WARN, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: ERROR, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: FATAL, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	time.Sleep(time.Second)
}

func TestRotateFileLog(t *testing.T) {
	logger := new(FileWriter)
	logger.SetFileBaseName("/var/loglet/log_test/test.loglet")
	logger.SetRotateSize(1024 * 100)
	logger.SetFileReserveNum(10)
	for i := 0; i < 1000000; i++ {
		logger.WriteLog(&LogMsg{msgLevel: DEBUG, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: INFO, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: WARN, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: ERROR, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: FATAL, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
	}

	time.Sleep(time.Second)
}
