package loglet

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

/**
 * 文件日志书写器定义
 */
type FileWriter struct {
	fileName          string
	rotateDaily       bool
	rotateSize        int64
	logFile           *os.File
	fileRollerCounter int //日志文件滚动计数器
	logFileReserveNum int //保留历史文件的个数
}

/**
 * 设置日志文件的基名
 */
func (logger *FileWriter) SetFileBaseName(fileName string) {
	logger.fileName = fileName
}

/**
 * 设置日志文件保留的个数
 */
func (logger *FileWriter) SetFileReserveNum(num int) {
	if num > 0 && num < 1000 {
		logger.logFileReserveNum = num
	} else {
		logger.logFileReserveNum = 10
	}
}

/**
 * 设置日志文件滚动的Size
 */
func (logger *FileWriter) SetRotateSize(size int64) {
	logger.rotateSize = size
}

/**
 * 向日志文件中输出日志
 */
func (logger *FileWriter) WriteLog(msg *LogMsg) {
	logger.fileRollerCounter++
	//为了避免频繁判断日志文件大小，导致性能下降，每写入1K条日志才判断是否要滚日志文件
	if logger.fileRollerCounter > 1000 {
		logger.rollLogFile()
		logger.deleteExpiredLogFile()
		logger.fileRollerCounter = 0
	}
	logFile, err := logger.getLoggingFile()
	if err != nil {
		printError("can not init log file: %s. error: %s.", logger.fileName, err.Error())
		return
	}
	_, err = logFile.WriteString(msg.getFormattedMsg())
	if err != nil {
		printError("can not write log to file: %s. error: %s.", logger.fileName, err.Error())
		return
	}
	//根据系统不同输入换行符
	switch runtime.GOOS {
	case "linux":
		logFile.WriteString("\n")
	case "windows":
		logFile.WriteString("\r\n")
	case "darwin":
		logFile.WriteString("\n")
	default:
		logFile.WriteString("\n")
	}
}

/**
 * 获取当前日志文件的句柄，不存在就打开一个新日志文件
 */
func (logger *FileWriter) getLoggingFile() (*os.File, error) {
	if logger.logFile != nil {
		return logger.logFile, nil
	}
	var err error
	logger.logFile, err = os.OpenFile(logger.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return logger.logFile, nil
}

/**
 * 如果单个日志大小超过配置size，则新建一个新的日志文件来写入
 */
func (logger *FileWriter) rollLogFile() error {
	if logger.fileName == "" || logger.logFile == nil {
		return nil
	}
	fileInfo, err := os.Stat(logger.fileName)

	if err != nil {
		printError("can not get log file status: %s.", err.Error())
		if os.IsNotExist(err) {
			logger.logFile = nil //如果日志文件在写入过程中被人为删除，则促使生成下一文件
		}
		return err
	}
	//如果没有定义文件大小，则默认一个日志文件100M
	rotateSize := logger.rotateSize
	if rotateSize == 0 {
		rotateSize = 1024 * 1024 * 100
	}
	if fileInfo.Size() < logger.rotateSize {
		return nil
	}

	// 修改log文件命名 liwenqiao 2017-7-20
	fileExt := filepath.Ext(logger.fileName)
	fileNameWithoutExt := strings.TrimSuffix(logger.fileName, fileExt)
	newlogFileName := fileNameWithoutExt + "." + time.Now().Format("20060102_150405") + fileExt
	err = logger.RenameCurLogFile(newlogFileName)
	if err != nil {
		printError("can not rename log file : %s.", err.Error())
		return err
	}

	return nil
}

/**
 * 清理过期日志
 */
func (logger *FileWriter) deleteExpiredLogFile() error {
	tempFileInfos, err := ioutil.ReadDir(filepath.Dir(logger.fileName))
	if err != nil {
		printError("can not get file infos : %s.", err.Error())
		return err
	}
	// 从所有文件中选择出.log文件，保存到fileInfos中
	fileInfos := make(SortableFileArr, 0, len(tempFileInfos))
	fileExt := filepath.Ext(logger.fileName) // fileExt == ".log"
	fileNameItems := strings.Split(logger.fileName, "/")
	fileName := fileNameItems[len(fileNameItems)-1]
	fileName = strings.Replace(fileName, fileExt, "", -1) + "."
	for i := 0; i < len(tempFileInfos); i++ {
		if !strings.HasPrefix(tempFileInfos[i].Name(), fileName) {
			continue
		}
		if !strings.HasSuffix(tempFileInfos[i].Name(), fileExt) {
			continue
		}
		fileInfos = append(fileInfos, tempFileInfos[i])
	}
	if len(fileInfos) == 1 {
		return nil //如果只有一个文件，则不删除，避免当前日志文件正处于写入状态
	}
	// 小于reserveNum时，不需要清理
	if len(fileInfos) <= logger.logFileReserveNum {
		return nil
	}
	sort.Sort(fileInfos)
	fileDir, _ := filepath.Split(logger.fileName)
	for i := 0; i < len(fileInfos)-logger.logFileReserveNum; i++ { //TODO: 逻辑是不是反了？
		fileName := fileDir + fileInfos[i].Name()
		err = os.Remove(fileName)
		if err != nil {
			printError("can not remove file : <%s> %s.", fileName, err.Error())
		}
	}
	return nil
}

/**
 * 重命名当前日志文件（例如在滚动日志文件时）
 */
func (logger *FileWriter) RenameCurLogFile(newFileName string) error {
	logger.Close()
	err := os.Rename(logger.fileName, newFileName)
	if err != nil {
		printError("can not rename log file: %s.", err.Error())
		return err
	}
	return nil
}

/**
 * 关闭当前日志文件
 */
func (logger *FileWriter) Close() {
	if logger.logFile != nil {
		err := logger.logFile.Close()
		if err != nil {
			printError("can not close log file: %s.", err.Error())
		}
		logger.logFile = nil
	}
}

/**
 * 声明一个排序数组，用于对文件名排序
 */
type SortableFileArr []os.FileInfo

func (s SortableFileArr) Less(i, j int) bool { return s[i].ModTime().Before(s[j].ModTime()) }

func (s SortableFileArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortableFileArr) Len() int { return len(s) }
