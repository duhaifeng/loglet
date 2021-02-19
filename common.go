package loglet

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

/**
 * 日志输出等级定义
 * @author duhaifeng
 */
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
	//定义一组数字映射等级，用于方便判断
	DEBUG_LEVEL = 0
	INFO_LEVEL  = 1
	WARN_LEVEL  = 2
	ERROR_LEVEL = 3
	FATAL_LEVEL = 4
)

/**
 * 日志输出消息包装
 * @author duhaifeng
 */
type LogMsg struct {
	msgLevel    string
	msgTime     time.Time
	targetPoint string
	msgContent  string
}

/**
 * 获取格式化后的日志输出字符
 * 2012-09-20 15:56:12  [ com.homer.HMain.printLog(HMain.java:24):java.lang.Class:http-bio-9980-exec-3:0 ] - [ DEBUG ]  log4j debug
 * @author duhaifeng
 */
func (msg *LogMsg) getFormattedMsg() string {
	timeStr := msg.msgTime.Format("2006-01-02 15:04:05.000")
	return timeStr + " " + msg.targetPoint + " [" + msg.msgLevel + "] " + msg.msgContent
}

/**
 * 从运行堆栈中获取日志产生的代码点
 * @author duhaifeng
 */
func getLoggingPoint() string {
	/**
	runtime.Stack()返回格式：
	goroutine 18 [running]:
	runtime/debug.Stack(0x0, 0x0, 0x0)
		/usr/local/go/src/runtime/debug/stack.go:24 +0xbe
	*/
	//查找当前的goroutine号（位于调用栈的的第一行中）
	stackBuf := make([]byte, 64)
	bufSize := runtime.Stack(stackBuf, false)
	stackBuf = stackBuf[:bufSize]
	stackBuf = bytes.TrimPrefix(stackBuf, []byte("goroutine "))
	stackBuf = stackBuf[:bytes.IndexByte(stackBuf, ' ')]
	routineNo := string(stackBuf)

	funcName := ""
	filePath := ""
	line := -1
	pcs := make([]uintptr, 25)
	callDepth := runtime.Callers(0, pcs)
	for i := 0; i < callDepth; i++ {
		pc := pcs[i]
		funcInfo := runtime.FuncForPC(pc)
		if strings.Contains(funcInfo.Name(), "loglet") {
			filePath, line = funcInfo.FileLine(pc)
			funcName = funcInfo.Name()
		}
	}
	file := filePath[strings.LastIndex(filePath, "/"):]
	return fmt.Sprintf("%s %d %s() [%s]", file, line, funcName, routineNo)
}

/**
 * 将日志内部错误打印到控制台
 * @author duhaifeng
 */
func printError(content string, contentArgs ...interface{}) {
	fmt.Fprintf(os.Stderr, "<loglet.error> "+content+"\n", contentArgs...)
}
