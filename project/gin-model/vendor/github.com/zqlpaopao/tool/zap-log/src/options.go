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

//InfoPathFileName string
//	warnPathFileName string
//	WithMaxAge int //*time.Hour
//	WithRotationCount uint
//	WithRotationTime int

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


