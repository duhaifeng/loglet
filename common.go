package loglet

import (
	"fmt"
	"os"
	"runtime/debug"
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
	callStack := string(debug.Stack())
	callStackArr := strings.Split(callStack, "\n")
	/**
	callStack返回格式：
	goroutine 18 [running]:
	runtime/debug.Stack(0x0, 0x0, 0x0)
		/usr/local/go/src/runtime/debug/stack.go:24 +0xbe
	loglet.getLoggingPoint(0x0, 0x0)
		/Users/duhf/Documents/IdeaProjects/goProjects/src/loglet/common.go:32 +0x4f
	loglet.(*loggerBase).getMsg(0xc420076438, 0x118e900, 0x11, 0xc420042e28, 0x3, 0x3, 0x0)
		/Users/duhf/Documents/IdeaProjects/goProjects/src/loglet/logger_base.go:27 +0x57
	loglet.(*loggerBase).Debug(0xc420076438, 0x118e900, 0x11, 0xc420042e28, 0x3, 0x3)
		/Users/duhf/Documents/IdeaProjects/goProjects/src/loglet/logger_base.go:31 +0x61
	loglet.TestLog(0xc4200b80f0)  （调用函数名字）
		/Users/duhf/Documents/IdeaProjects/goProjects/src/loglet/logger_test.go:27 +0x40d （文件行号位置）
	testing.tRunner(0xc4200b80f0, 0x11950a8)
		/usr/local/go/src/testing/testing.go:746 +0x11f
	created by testing.(*T).Run
		/usr/local/go/src/testing/testing.go:789 +0x4e3
	*/
	//查找当前的goroutine号
	routineNo := "-1"
	for i := 0; i < len(callStackArr); i++ {
		line := callStackArr[i]
		if strings.HasPrefix(line, "goroutine") {
			lineItems := strings.Split(line, " ")
			if len(lineItems) > 1 {
				routineNo = lineItems[1]
			}
		}
	}

	loggingFileLine := ""
	loggingFunc := ""
	lastLine := ""
	lastTwoLine := ""
	//查找上层代码记录日志的调用点
	for i := len(callStackArr) - 1; i >= 0; i-- {
		line := callStackArr[i]
		//如果进入了common/loglet.go或者logger_base.go，就说明上一级调用是外层业务代码的Log调用点，then this is the loglet point.
		if loggingFileLine == "" && (strings.Contains(line, "common/loglet.go") || strings.Contains(line, "logger_base.go")) {
			lastLine = strings.Split(lastLine, "(0x")[0]
			funcItems := strings.Split(lastLine, "/")
			//funcItems = strings.Split(funcItems[len(funcItems)-1], ".")
			loggingFunc = funcItems[len(funcItems)-1]

			//原始格式为：“/home/jupiter/cmd/joypaw-cli/xyz.go:35 +0xa4”，需要提炼出其中的 xyz.go:35
			lastTwoLine = strings.Split(lastTwoLine, " ")[0]
			lineItems := strings.Split(lastTwoLine, "/")
			loggingFileLine = lineItems[len(lineItems)-1]
			break
		}
		lastTwoLine = lastLine
		lastLine = line
	}
	return loggingFileLine + " " + loggingFunc + "() [" + routineNo + "]"
}

/**
 * 将日志内部错误打印到控制台
 * @author duhaifeng
 */
func printError(content string, contentArgs ...interface{}) {
	fmt.Fprintf(os.Stderr, "<loglet.error> "+content+"\n", contentArgs...)
}
