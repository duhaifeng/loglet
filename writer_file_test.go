package loglet

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const line = "ABCDEFGHIJKLMNOPQRSTUVWXYZ---日志--- "

func TestFileWriter(t *testing.T) {
	logger := new(FileWriter)
	logger.Init()
	logger.SetFileBaseName("/var/log/log_test/test.log")
	logger.WriteLog(&LogMsg{msgLevel: DEBUG, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: INFO, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: WARN, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: ERROR, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	logger.WriteLog(&LogMsg{msgLevel: FATAL, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line})
	time.Sleep(time.Second)
}

func TestRotateFileLog(t *testing.T) {
	logger := new(FileWriter)
	logger.Init()
	logger.SetFileBaseName("D:/log_test/test.log")
	logger.SetRotateSize(1024 * 1000 * 10)
	logger.SetFileReserveNum(10)
	lastSecond := time.Now().Format("15:04:05")
	var lastSecondWriteCount int32 = 0
	for i := 0; i < 1000000; i++ {
		logger.WriteLog(&LogMsg{msgLevel: DEBUG, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: INFO, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: WARN, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: ERROR, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		logger.WriteLog(&LogMsg{msgLevel: FATAL, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: line + strconv.Itoa(i)})
		if lastSecond == time.Now().Format("15:04:05") {
			atomic.AddInt32(&lastSecondWriteCount, 5)
		} else {
			fmt.Printf("write %d log at %s\n", lastSecondWriteCount, lastSecond)
			lastSecond = time.Now().Format("15:04:05")
			atomic.StoreInt32(&lastSecondWriteCount, 0)
		}
	}
	time.Sleep(time.Second)
}

func TestMultiRoutineFileLog(t *testing.T) {
	logger := new(FileWriter)
	logger.Init()
	logger.SetFileBaseName("D:/log_test/test.log")
	logger.SetRotateSize(1024 * 1000 * 10)
	logger.SetFileReserveNum(10)
	lastSecond := time.Now().Format("15:04:05")
	var lastSecondWriteCount int32 = 0

	var wg sync.WaitGroup
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func(worker string) {
			for i := 0; i < 1000000; i++ {
				logger.WriteLog(&LogMsg{msgLevel: DEBUG, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: worker + line + strconv.Itoa(i)})
				logger.WriteLog(&LogMsg{msgLevel: INFO, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: worker + line + strconv.Itoa(i)})
				logger.WriteLog(&LogMsg{msgLevel: WARN, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: worker + line + strconv.Itoa(i)})
				logger.WriteLog(&LogMsg{msgLevel: ERROR, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: worker + line + strconv.Itoa(i)})
				logger.WriteLog(&LogMsg{msgLevel: FATAL, msgTime: time.Now(), targetPoint: getLoggingPoint(0), msgContent: worker + line + strconv.Itoa(i)})
				if lastSecond == time.Now().Format("15:04:05") {
					atomic.AddInt32(&lastSecondWriteCount, 5)
				} else {
					fmt.Printf("write %d log at %s\n", lastSecondWriteCount, lastSecond)
					lastSecond = time.Now().Format("15:04:05")
					atomic.StoreInt32(&lastSecondWriteCount, 0)
				}
			}
			wg.Done()
		}(fmt.Sprintf("worker-%d ", w))
	}
	wg.Wait()
	time.Sleep(time.Second)
}

func TestFilePath(t *testing.T) {
	fmt.Println(filepath.Dir("/var/log/log_test/test.log"))
	fmt.Println(filepath.Dir("D:/log_test/test.log"))
	logger := new(FileWriter)
	fmt.Println(logger.isPathExists("D:/test.log"))
	fmt.Println(logger.isPathExists("D:/log_test/"))
	fmt.Println(logger.isPathExists("D:/log_test/test.log"))
}
