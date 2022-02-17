package src

var logAsync = new(syncLogConfig)

type syncLogConfig struct {
	poolHandler *pool
	//syncTime int
	buffSize int
	//syncNum int
	syncGoNum int
}

// DebugAsync level
func DebugAsync(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg, debugLevel, args)
}

// InfoAsync level
func InfoAsync(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg, infoLevel, args)
}

// WarnAsync level
func WarnAsync(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg, warnLevel, args)
}

//ErrorAsync level
func ErrorAsync(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg, errorLevel, args)
}

//MsgAsync Really write
func (e *ErrorHandle) MsgAsync(err string) {
	task := InitTask(e.Msg, err)
	logAsync.poolHandler.AddTask(task)
}


//initSyncGoPool Initialize synchronization process pool
func (log *syncLogConfig) initSyncGoPool() {
	//创建协程池
	log.poolHandler = NewPool(log.syncGoNum, log.buffSize)
	//启动协程池
	go log.poolHandler.Run()
	logAsync = log
}
