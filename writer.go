package loglet

import (
	"fmt"
	"os"
)

/**
 * 日志书写器抽象定义
 */
type LogWriter interface {
	WriteLog(msg *LogMsg)
	Close()
}

/**
 * 控制台日志书写器定义
 */
type ConsoleWriter struct {
}

/**
 * 向控制台输出日志
 */
func (logger *ConsoleWriter) WriteLog(msg *LogMsg) {
	if msg.msgLevel == ERROR || msg.msgLevel == FATAL {
		fmt.Fprintf(os.Stderr, msg.getFormattedMsg()+"\n")
	} else {
		fmt.Println(msg.getFormattedMsg())
	}
}

/**
 * 关闭控制台日志书写器（控制台本身不需要关闭，为了实现多态，这里补足Close方法）
 */
func (logger *ConsoleWriter) Close() {
}
