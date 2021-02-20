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
 */
func (msg *LogMsg) getFormattedMsg() string {
	timeStr := msg.msgTime.Format("2006-01-02 15:04:05.000")
	return timeStr + " " + msg.targetPoint + " [" + msg.msgLevel + "] " + msg.msgContent
}

/**
 * 从运行堆栈中获取日志产生的代码点
 */
func getLoggingPoint(offset int) string {
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

	//查找日志记录点的函数名、文件及行号
	outerCallerIndex := 0
	pcs := make([]uintptr, 25)
	callDepth := runtime.Callers(0, pcs)
	for i := callDepth - 1; i >= 0; i-- {
		pc := pcs[i]
		funcInfo := runtime.FuncForPC(pc)
		//如果发现了loglet包路径，则说明进入了日志模块内部，获取上一次调用作为日志记录点
		if strings.Contains(funcInfo.Name(), "loglet") {
			outerCallerIndex = i + 1
			break
		}
	}
	//将日志记录点进行偏移（前提是外部进行了指定）
	if outerCallerIndex+offset < callDepth && outerCallerIndex+offset >= 0 {
		outerCallerIndex += offset
	}
	outerCallerPc := pcs[outerCallerIndex]
	funcInfo := runtime.FuncForPC(outerCallerPc)
	filePath, line := funcInfo.FileLine(outerCallerPc)
	funcName := funcInfo.Name()
	file := filePath[strings.LastIndex(filePath, "/")+1:]
	return fmt.Sprintf("%s %d %s() [%s]", file, line, funcName, routineNo)
}

/**
 * 将日志内部错误打印到控制台
 */
func printError(content string, contentArgs ...interface{}) {
	fmt.Fprintf(os.Stderr, "<log.error> "+content+"\n", contentArgs...)
}
