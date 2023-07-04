package src

var logAsync = new(syncLogConfig)

type syncLogConfig struct {
	poolHandler *Pool
	//syncTime int
	buffSize int
	//syncNum int
	syncGoNum int
}

// DebugAsync level
func DebugAsync(tag string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(tag, debugLevel, args)
}

// InfoAsync level
func InfoAsync(tag string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(tag, infoLevel, args)
}

// WarnAsync level
func WarnAsync(tag string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(tag, warnLevel, args)
}

// ErrorAsync level
func ErrorAsync(tag string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(tag, errorLevel, args)
}

// MsgAsync Really write
func (e *ErrorHandle) MsgAsync(msg string) {
	task := InitTask(e.Msg, msg)
	logAsync.poolHandler.AddTask(task)
}

// initSyncGoPool Initialize synchronization process pool
func (log *syncLogConfig) initSyncGoPool() {
	//创建协程池
	log.poolHandler = NewPool(log.syncGoNum, log.buffSize)
	//启动协程池
	go log.poolHandler.Run()
	logAsync = log
}
