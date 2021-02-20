package loglet

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logger := NewLogger()
	conf := make(map[string]string)
	conf["writers"] = "file,console"
	conf["log_file"] = "/var/log/log_test/test.log"
	conf["log_level"] = "info"
	conf["max_size"] = "150k"
	conf["file_number"] = "10"
	logger.Init(conf)
	for i := 0; i < 100; i++ {
		logger.Debug("d:%d, f:%f, s:%s ", 111, 222.2, line)
		logger.Info("d:%d, f:%f, s:%s ", 111, 222.2, line)
		logger.Warn("d:%d, f:%f, s:%s ", 111, 222.2, line)
	}

	time.Sleep(time.Second)
}

func TestConcurrentLog(t *testing.T) {
	logger := NewLogger()
	conf := make(map[string]string)
	conf["writers"] = "console"
	conf["log_file"] = "/var/log/log_test/test.log"
	conf["log_level"] = "debug"
	conf["max_size"] = "150k"
	conf["file_number"] = "10"
	logger.Init(conf)
	for n := 0; n < 10; n++ {
		go func() {
			for i := 0; i < 10000; i++ {
				logger.Debug("%s ", line)
			}
		}()
	}

	time.Sleep(time.Second * 10)
}
