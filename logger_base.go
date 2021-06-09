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
	logPositionOffset int                  //允许外部定义一个偏移量，避免外部二次封装时日志都打在外面的封装点上
	logWriters        map[string]LogWriter //为了防止配置中重复出现file、console等，采用map进行滤重
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
 * 从日志缓存管道中获取日志记录
 */
func (logger *loggerBase) RegisterWriter(name string, logWriter LogWriter) {
	if logger.logWriters == nil {
		logger.logWriters = make(map[string]LogWriter)
	}
	logger.logWriters[name] = logWriter
}

/**
 * 关闭所有的日志书写器
 */
func (logger *loggerBase) CloseWriters() {
	for _, logWriter := range logger.logWriters {
		logWriter.Close()
	}
	logger.logWriters = nil
}

/**
 * 向日志缓存管道缓存日志
 */
func (logger *loggerBase) writeLog(msg *LogMsg) {
	for _, logWriter := range logger.logWriters {
		logWriter.WriteLog(msg)
	}
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
	logger.writeLog(msg)
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
	logger.writeLog(msg)
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
	logger.writeLog(msg)
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
	logger.writeLog(msg)
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
	logger.writeLog(msg)
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
