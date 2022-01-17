package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"time"
)
const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)
type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TracePlugin{}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	//_ctx := db.Statement.Context
	//ctx, ok := _ctx.(core.Context)
	//if !ok {
	//	return
	//}

	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	fmt.Println("执行的sql语句",sql)
	fmt.Println("执行的时常",time.Since(ts).Seconds())
	fmt.Println("文件及行号",utils.FileWithLineNum())
	fmt.Println("DB影响的行数",db.Statement.RowsAffected)

	//sqlInfo := new(trace.SQL)
	//sqlInfo.Timestamp = time_parse.CSTLayoutString()
	//sqlInfo.SQL = sql
	//sqlInfo.Stack = utils.FileWithLineNum()
	//sqlInfo.Rows = db.Statement.RowsAffected
	//sqlInfo.CostSeconds = time.Since(ts).Seconds()
	//ctx.Trace().AppendSQL(sqlInfo)

	return
}


