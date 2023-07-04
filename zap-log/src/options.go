package src

import "errors"

// An Option configures a Logger.
type Option interface {
	apply(*LogConfig)
}

// OptionFunc wraps a func  it satisfies the Option interface.
type OptionFunc func(*LogConfig)
type CallBackFunc func(int, string, *CallBack)

func (f OptionFunc) apply(log *LogConfig) {
	f(log)
}

// NewLogConfig init log
func NewLogConfig(f ...Option) *LogConfig {
	log := &LogConfig{}
	return log.WithOptions(f...)
}

// clone is copy
func (l *LogConfig) clone() *LogConfig {
	cy := *l
	return &cy
}

// WithOptions Why is there no address? In this way, multiple l can be copied without conflict
func (l *LogConfig) WithOptions(opts ...Option) *LogConfig {
	c := l.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// InitInfoPathFileName Initialize the log path of info debug
func InitInfoPathFileName(path string) OptionFunc {
	return func(l *LogConfig) {
		l.infoPathFileName = path
	}
}

// InitWarnPathFileName Initialize the log path of warn error
func InitWarnPathFileName(path string) OptionFunc {
	return func(l *LogConfig) {
		l.warnPathFileName = path
	}
}

// InitWithMaxAge Initialize log save time
func InitWithMaxAge(i int) OptionFunc {
	return func(l *LogConfig) {
		l.withMaxAge = i
	}
}

// InitWithRotationCount Initialize log save time
func InitWithRotationCount(i uint) OptionFunc {
	return func(l *LogConfig) {
		l.withRotationCount = i
	}
}

// InitWithRotationTime Initialize log rotation cycle
func InitWithRotationTime(i int) OptionFunc {
	return func(l *LogConfig) {
		l.withRotationTime = i
	}
}

// InitWithIp Whether IP is recorded during initialization
func InitWithIp(i int8) OptionFunc {
	return func(l *LogConfig) {
		l.ipTag = i
	}
}

// InitBufferSize Initialize the size of the resource pool
func InitBufferSize(i int) OptionFunc {
	if err := checkSize(i, bufferMax); nil != err {
		panic(err)
	}
	return func(l *LogConfig) {
		NewBuffErrSize(i)
		NewTaskBufferPool(i)
		NewCallBackBuffSize(i)
	}
}

// InitCallFunc Initialize the size of the resource pool
func InitCallFunc(f CallBackFunc) OptionFunc {
	if f == nil {
		return func(l *LogConfig) {
			NewCallBackFunc()
		}
	}
	return func(l *LogConfig) {
		callBackFunc = f
	}
}

/*************************************************** syncLogConfig ***************************************************/

type OptionAsync interface {
	apply(*syncLogConfig)
}
type OptionSyncFunc func(*syncLogConfig)

func (f OptionSyncFunc) apply(log *syncLogConfig) {
	f(log)
}

// NewAsyncLogConfig init log
func NewAsyncLogConfig(f ...OptionAsync) {
	logAsync = &syncLogConfig{
		buffSize:    syncBuffSize,
		syncGoNum:   syncGoNum,
		poolHandler: &Pool{},
	}
	logAsync.WithOptions(f...).initSyncGoPool()
}

// clone is copy
func (log *syncLogConfig) clone() *syncLogConfig {
	cy := *log
	return &cy
}

// WithOptions Why is there no address? In this way, multiple l can be copied without conflict
func (log *syncLogConfig) WithOptions(opts ...OptionAsync) *syncLogConfig {
	c := log.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// InitLogAsyncBuffSize Initialize log save time
func InitLogAsyncBuffSize(i int) OptionSyncFunc {
	if err := checkSize(i, maxSyncBuffSize); nil != err {
		panic(err)
	}
	return func(l *syncLogConfig) {
		l.buffSize = i
	}
}

// checkSize check size
func checkSize(i, max int) error {
	if i > max {
		return errors.New("i is more than max")
	}
	return nil
}

// InitLogAsyncGoNum Initialize log go num
func InitLogAsyncGoNum(i int) OptionSyncFunc {
	return func(l *syncLogConfig) {
		l.syncGoNum = i
	}
}
