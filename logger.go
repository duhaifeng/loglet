package loglet

import (
	"strconv"
	"strings"
)

/**
 * 创建一个日志实例
 */
func NewLogger() *Logger {
	logger := new(Logger)
	defaultConfs := make(map[string]string)
	defaultConfs["writers"] = "console"
	defaultConfs["log_level"] = "debug"
	logger.Init(defaultConfs)
	return logger
}

/**
 * 日志实例对象封装
 */
type Logger struct {
	loggerBase
}

/**
 * 初试化日志实例配置，如果不传入任何配置，则只向控制台输出
 */
func (logger *Logger) Init(configs map[string]string) {
	//为了避免重复Init，需要先关闭现已打开资源
	logger.CloseWriters()

	logger.loggerBase.SetLogLevel(configs["log_level"])
	logger.logWriters = make(map[string]LogWriter)
	writers := strings.Split(configs["writers"], ",")
	for _, writerName := range writers {
		if strings.TrimSpace(writerName) == "console" {
			logger.RegisterWriter("console", logger.createConsoleWriter(configs))
		}
		if strings.TrimSpace(writerName) == "file" {
			logger.RegisterWriter("file", logger.createFileWriter(configs))
		}
	}
	if len(logger.logWriters) == 0 {
		logger.RegisterWriter("console", logger.createConsoleWriter(configs))
	}
}

/**
 * 创建一个控制台日志书写器
 */
func (logger *Logger) createConsoleWriter(configs map[string]string) *ConsoleWriter {
	consoleLogger := new(ConsoleWriter)
	return consoleLogger
}

/**
 * 创建一个文件日志书写器
 */
func (logger *Logger) createFileWriter(configs map[string]string) *FileWriter {
	fileLogger := new(FileWriter)
	fileLogger.Init()
	fileLogger.SetFileBaseName(configs["log_file"])
	fileSizeStr := strings.ToUpper(configs["max_size"])
	fileSizeUnit := fileSizeStr[len(fileSizeStr)-1:] //取配置的最后一个字母作为日志文件大小的单位
	fileSizeStr = strings.Replace(fileSizeStr, "K", "", -1)
	fileSizeStr = strings.Replace(fileSizeStr, "M", "", -1)
	fileSizeStr = strings.Replace(fileSizeStr, "G", "", -1)
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		printError("log file size config error. use default size: 100M")
		fileSize = 1024 * 1024 * 100
	} else {
		// modified by liwenqiao 2017-7-20
		if fileSizeUnit == "K" {
			fileSize = fileSize * 1024
		} else if fileSizeUnit == "G" {
			fileSize = fileSize * 1024 * 1024 * 1024
		} else {
			fileSize = fileSize * 1024 * 1024
		}
	}
	fileLogger.SetRotateSize(int64(fileSize))
	//设置要保留日志文件的个数
	fileNum, err := strconv.Atoi(configs["file_number"])
	if err != nil {
		printError("log file reserve number config error. use default: 10")
		fileLogger.SetFileReserveNum(10)
	} else {
		fileLogger.SetFileReserveNum(fileNum)
	}
	return fileLogger
}
