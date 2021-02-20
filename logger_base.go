package loglet

import (
	"fmt"
	"strings"
	"time"
)

/**
 * 日志记录器对象抽象定义
 */
type loggerBase struct {
	logLevel          int
	logPositionOffset int //允许外部定义一个偏移量，避免外部二次封装时日志都打在外面的封装点上
	msgChannel        chan *LogMsg
}

/**
 * 初始化日志记录器内部日志缓存管道
 */
func (logger *loggerBase) Init() {
	logger.msgChannel = make(chan *LogMsg, 100000)
}

/**
 * 设置日志输出级别
 */
func (logger *loggerBase) SetLogLevel(level string) {
	logger.logLevel = logger.getLogLevelNum(level)
}

/**
 * 设置日志打印堆栈点的偏移量
 */
func (logger *loggerBase) SetLogPositionOffset(logPositionOffset int) {
	logger.logPositionOffset = logPositionOffset
}

/**
 * 向日志缓存管道中追加日志记录
 */
func (logger *loggerBase) AppendMsg(msg *LogMsg) {
	if logger.msgChannel != nil {
		logger.msgChannel <- msg
	}
}

/**
 * 从日志缓存管道中获取日志记录
 */
func (logger *loggerBase) WaitMsg() *LogMsg {
	return <-logger.msgChannel
}

/**
 * 将一条上次传入的消息进行封装
 */
func (logger *loggerBase) getMsg(msg string, msgArgs ...interface{}) *LogMsg {
	return &LogMsg{msgTime: time.Now(), targetPoint: getLoggingPoint(logger.logPositionOffset), msgContent: fmt.Sprintf(msg, msgArgs...)}
}

/**
 * 写入Debug级别日志
 */
func (logger *loggerBase) Debug(content string, contentArgs ...interface{}) {
	if logger.logLevel > DEBUG_LEVEL {
		return
	}
	msg := logger.getMsg(content, contentArgs...)
	msg.msgLevel = DEBUG
	logger.sendMsg(msg)
}

/**
 * 写入Info级别日志
 */
func (logger *loggerBase) Info(content string, contentArgs ...interface{}) {
	if logger.logLevel > INFO_LEVEL {
		return
	}
	msg := logger.getMsg(content, contentArgs...)
	msg.msgLevel = INFO
	logger.sendMsg(msg)
}

/**
 * 写入Warning级别日志
 */
func (logger *loggerBase) Warn(content string, contentArgs ...interface{}) {
	if logger.logLevel > WARN_LEVEL {
		return
	}
	msg := logger.getMsg(content, contentArgs...)
	msg.msgLevel = WARN
	logger.sendMsg(msg)
}

/**
 * 写入Error级别日志
 */
func (logger *loggerBase) Error(content interface{}, contentArgs ...interface{}) {
	if logger.logLevel > ERROR_LEVEL {
		return
	}
	var msg *LogMsg
	errContent, ok := content.(error)
	if ok {
		msg = logger.getMsg(errContent.Error(), contentArgs...)
	}
	contentStr, ok := content.(string)
	{
		msg = logger.getMsg(contentStr, contentArgs...)
	}

	msg.msgLevel = ERROR
	logger.sendMsg(msg)
}

/**
 * 写入Fatal级别日志
 */
func (logger *loggerBase) Fatal(content string, contentArgs ...interface{}) {
	if logger.logLevel > FATAL_LEVEL {
		return
	}
	msg := logger.getMsg(content, contentArgs...)
	msg.msgLevel = FATAL
	logger.sendMsg(msg)
}

/**
 * 向日志缓存管道缓存日志
 */
func (logger *loggerBase) sendMsg(msg *LogMsg) {
	if logger.msgChannel == nil {
		return
	}
	logger.msgChannel <- msg
}

/**
 * 关闭日志缓存管道
 */
func (logger *loggerBase) CloseChannel() {
	if logger.msgChannel == nil {
		return
	}
	close(logger.msgChannel)
	logger.msgChannel = nil
}

/**
 * 判断当前日志记录是否达到了输出的定义级别，如果未达到则丢弃上层传入的消息
 */
func (logger *loggerBase) matchLogLevel(msg *LogMsg) bool {
	return logger.getLogLevelNum(msg.msgLevel) >= logger.logLevel
}

/**
 * 获取日志级别对应的数字，便于判断
 */
func (logger *loggerBase) getLogLevelNum(level string) int {
	level = strings.ToUpper(level)
	switch level {
	case DEBUG:
		return DEBUG_LEVEL
	case INFO:
		return INFO_LEVEL
	case WARN:
		return WARN_LEVEL
	case ERROR:
		return ERROR_LEVEL
	case FATAL:
		return FATAL_LEVEL
	default:
		return -1
	}
}
