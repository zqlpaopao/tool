package src

// An Option configures a Logger.
type Option interface {
	apply(*logConfig)
}
//OptionFunc wraps a func  it satisfies the Option interface.
type OptionFunc func(*logConfig)

func (f OptionFunc) apply(log *logConfig) {
	f(log)
}

//NewLogConfig init log
func NewLogConfig(f... Option)*logConfig{
	log := &logConfig{}
	return log.WithOptions(f...)
}

//clone is copy
func (l *logConfig) clone() *logConfig {
	cy := *l
	return &cy
}

//WithOptions Why is there no address? In this way, multiple l can be copied without conflict
func (l logConfig)WithOptions(opts ...Option) *logConfig {
	c := l.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}


//InitInfoPathFileName Initialize the log path of info debug
func InitInfoPathFileName(path string)OptionFunc{
	return func(l *logConfig){
		l.infoPathFileName = path
	}
}

//InitWarnPathFileName Initialize the log path of warn error
func InitWarnPathFileName(path string)OptionFunc{
	return func(l *logConfig) {
		l.warnPathFileName = path
	}
}

//InitWithMaxAge Initialize log save time
func InitWithMaxAge(i int)OptionFunc{
	return func(l *logConfig) {
		l.withMaxAge = i
	}
}

//InitWithRotationCount Initialize log save time
func InitWithRotationCount(i uint)OptionFunc{
	return func(l *logConfig) {
		l.withRotationCount = i
	}
}

//InitWithRotationTime Initialize log rotation cycle
func InitWithRotationTime(i int)OptionFunc{
	return func(l *logConfig) {
		l.withRotationTime = i
	}
}

//InitWithIp Whether IP is recorded during initialization
func InitWithIp(i int8)OptionFunc{
	return func(l *logConfig) {
		l.ipTag = i
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

//NewAsyncLogConfig init log
func NewAsyncLogConfig(f... OptionAsync){
	logAsync = &syncLogConfig{
		buffSize:  syncBuffSize,
		syncGoNum: syncGoNum,
		poolHandler :&pool{},
	}
	logAsync.WithOptions(f...).initSyncGoPool()
}

//clone is copy
func (log *syncLogConfig) clone() *syncLogConfig {
	cy := *log
	return &cy
}

//WithOptions Why is there no address? In this way, multiple l can be copied without conflict
func (log syncLogConfig)WithOptions(opts ...OptionAsync) *syncLogConfig {
	c := log.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

//InitLogAsyncBuffSize Initialize log save time
func InitLogAsyncBuffSize(i int)OptionSyncFunc{
	return func(l *syncLogConfig) {
		l.buffSize = i
	}
}

//InitLogAsyncGoNum Initialize log go num
func InitLogAsyncGoNum(i int)OptionSyncFunc{
	return func(l *syncLogConfig) {
		l.syncGoNum = i
	}
}
