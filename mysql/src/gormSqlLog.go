package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/zqlpaopao/tool/format/src"
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"reflect"
	"strings"
)

type GormSqlLog struct {
	sqlString string
	printFlag bool
}

var gormSqlLogMan = new(GormSqlLog)

func GetGormSqlLogMan(dbs *gorm.DB, flag bool) {
	gormSqlLogMan.sqlString = ""
	gormSqlLogMan.printFlag = flag
	dbs.SetLogger(gormSqlLogMan)
	return
}

//Print Gorm implements the Print method, only need to have method inheritance implementation Print
func (g *GormSqlLog) Print(v ...interface{}) {
	if v[0] != "sql" {
		return
	}
	sqlString := v[3].(string)
	var list []interface{}
	if reflect.TypeOf(v[4]).Kind() == reflect.Slice {
		s := reflect.ValueOf(v[4])
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface())
		}
	}
	for _, v := range list {
		g.sqlString = strings.Replace(sqlString, "?", pkg.StringFromAssertionFloat(v), 1)
	}
	if g.printFlag {
		src.PrintGreen(g.sqlString)
	}
}

//String get sql Info
func (g *GormSqlLog) String() string {
	return g.sqlString
}
