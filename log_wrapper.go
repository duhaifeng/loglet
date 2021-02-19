package loglet

type LogWrapper struct {
	logger *Logger
}

func (this *LogWrapper) SetOriginalLogger(logger *Logger) {
	this.logger = logger
}

func (this *LogWrapper) GetOriginalLogger() *Logger {
	return this.logger
}

func (this *LogWrapper) appendReqIdToLog(reqId, content string) string {
	if reqId == "" {
		reqId = "<no-request-id>"
	}
	return reqId + " " + content
}

func (this *LogWrapper) Debug(reqId, content string, contentArgs ...interface{}) {
	this.logger.Debug(this.appendReqIdToLog(reqId, content), contentArgs...)
}

func (this *LogWrapper) Info(reqId, content string, contentArgs ...interface{}) {
	this.logger.Info(this.appendReqIdToLog(reqId, content), contentArgs...)
}

func (this *LogWrapper) Warn(reqId, content string, contentArgs ...interface{}) {
	this.logger.Warn(this.appendReqIdToLog(reqId, content), contentArgs...)
}

func (this *LogWrapper) Error(reqId, content string, contentArgs ...interface{}) {
	this.logger.Error(this.appendReqIdToLog(reqId, content), contentArgs...)
}
