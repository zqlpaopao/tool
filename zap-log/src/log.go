package src

import (
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	ipInfo "github.com/zqlpaopao/tool/ip/src"
	"go.uber.org/zap"
	zapCore "go.uber.org/zap/zapcore"
	"io"
	"time"
)

const (
	_ = iota
	debugLevel
	infoLevel
	warnLevel
	errorLevel
)

type logConfig struct {
	infoPathFileName  string
	warnPathFileName  string
	withMaxAge        int //*time.Hour
	withRotationCount uint
	withRotationTime  int
	ipTag             int8
}

var (
	Logger *zap.Logger
	IpInfo string
)


func init(){
	errHandler = new(ErrorHandle)
}


//InitLoggerHandler -- ----------------------------
//--> @Description  Initialize log processing assistant
//--> @Param
//--> @return
//-- ----------------------------
func InitLoggerHandler(logConf *logConfig) {
	checkLogConfig(logConf)
	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	//infoWriter := getWriter(logConf.infoPathFileName)
	//warnWriter := getWriter(logConf.warnPathFileName)
	infoLevel, warnLevel := checkLevel()
	// 最后创建具体的Logger
	Logger = zap.New(zapCore.NewTee(
		zapCore.NewCore(getEncoder(), zapCore.AddSync(getWriter(logConf.infoPathFileName, logConf)), infoLevel),
		zapCore.NewCore(getEncoder(), zapCore.AddSync(getWriter(logConf.warnPathFileName, logConf)), warnLevel),
	), zap.AddCaller(), zap.AddCallerSkip(1)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
}

//checkLogConfig -- ----------------------------
//--> @Description check Args
//--> @Param
//--> @return
//-- ----------------------------
func checkLogConfig(logConf *logConfig) {
	if logConf.warnPathFileName == "" || logConf.infoPathFileName == "" {
		panic("Empty directory is not allowed")
	}
}

//getEncoder -- ----------------------------
//--> @Description  Initialize configuration
//--> @Param
//--> @return
//-- ----------------------------
func getEncoder() zapCore.Encoder {
	//encoder := zapCore.NewConsoleEncoder(zapCore.EncoderConfig{//只有参数是json格式
	return zapCore.NewJSONEncoder(zapCore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapCore.CapitalLevelEncoder, //level转换为全大写
		//EncodeLevel:    zapCore.LowercaseLevelEncoder,//小写
		TimeKey: "time",
		EncodeTime: func(t time.Time, enc zapCore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapCore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapCore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
		LineEnding: zapCore.DefaultLineEnding,
	})
}

//checkLevel -- --------------------------------------
//--> @Description checkLevel info warn
//--> @Param
//--> @return
//-- ----------------------------
func checkLevel() (zap.LevelEnablerFunc, zap.LevelEnablerFunc) {
	// 实现两个判断日志等级的interface (其实 zapCore.*Level 自身就是 interface)
	//zap.LevelEnablerFunc()
	return func(lvl zapCore.Level) bool { //infoLevel
			return lvl < zapCore.WarnLevel
		},
		func(lvl zapCore.Level) bool { //WarnLevel
			return lvl >= zapCore.WarnLevel
		}
}

//getWriter -- ----------------------------
//--> @Description Get auto split function
//--> @Param
//--> @return
//-- ----------------------------
func getWriter(filename string, logConf *logConfig) io.Writer {
	// 生成rotateLogs的Logger 实际生成的文件名 demo.log.YY_mm_dd_HH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotateLogs.New(
		//%Y%m%d%H%M%S 年月日 时分秒 能记录到小时
		filename,
		getStartRotateLogsConf(logConf)...,
	)
	if err != nil {
		panic(err)
	}
	return hook
}

//-- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func getStartRotateLogsConf(logConf *logConfig) (op []rotateLogs.Option) {
	//保留最大文件数
	if logConf.withRotationCount > 0 {
		op = append(op, rotateLogs.WithRotationCount(logConf.withRotationCount))
	} else {
		op = append(op, rotateLogs.WithRotationCount(0)) //禁用

	}

	//文件保留时间
	if logConf.withMaxAge > 0 {
		op = append(op, rotateLogs.WithMaxAge(time.Hour*time.Duration(logConf.withMaxAge)))
	} else {
		op = append(op, rotateLogs.WithMaxAge(-1))
	}

	//文件轮换间隔。默认情况下，日志每 86400 秒轮换一次。注意：请记住使用 time.Duration 值。
	if logConf.withRotationTime < 0 {
		op = append(op, rotateLogs.WithRotationTime(time.Hour))
	}
	if logConf.withRotationTime > 0 {
		op = append(op, rotateLogs.WithRotationTime(time.Duration(logConf.withRotationTime)))
	}

	if logConf.ipTag > 0 {
		IpInfo = ipInfo.GetEth0()
	}
	return op
}

// FormatLog Format str
//func FormatLog(args []interface{}) *zap.Logger {
//	log := Logger.With(ToJsonData(args))
//	return log
//}

type ErrorHandle struct {
	msg  string
	args []interface{}
	tag  int
}

// Debug level
func Debug(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg,debugLevel,args...)
}

// Info level
func Info(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg,infoLevel,args...)
}

// Warn level
func Warn(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg,warnLevel,args...)
	//FormatLog(args).Sugar().Warnf(msg)
}

//Error level
func Error(msg string, args ...interface{}) *ErrorHandle {
	return errHandler.initParams(msg,errorLevel,args...)
}

// Msg Really write
func (e *ErrorHandle) Msg(err string) {
	switch e.tag {
	case debugLevel:
		Logger.With(ToJsonIpInfo(), ToJsonData(e.args), ToJsonError(err)).Sugar().Debug(e.msg)
	case infoLevel:
		Logger.With(ToJsonIpInfo(), ToJsonData(e.args), ToJsonError(err)).Sugar().Infof(e.msg)
	case warnLevel:
		Logger.With(ToJsonIpInfo(), ToJsonData(e.args), ToJsonError(err)).Sugar().Warn(e.msg)
	case errorLevel:
		Logger.With(ToJsonIpInfo(), ToJsonData(e.args), ToJsonError(err)).Sugar().Error(e.msg)
	}
}

// ToJsonData to string-byte
func ToJsonData(args []interface{}) zap.Field {
	det := make([]string, 0)
	if len(args) > 0 {
		for _, v := range args {
			det = append(det, fmt.Sprintf("%+v", v))
		}
	}
	z := zap.Any("params", det)
	return z
}

//ToJsonError error info
func ToJsonError(err string) zap.Field {
	z := zap.Any("errMsg", err)
	return z
}

//ToJsonIpInfo ipInfo
func ToJsonIpInfo() zap.Field {
	z := zap.Any("localIp", IpInfo)
	return z
}
