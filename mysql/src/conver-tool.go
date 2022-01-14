package mysql

import (
	"errors"
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"reflect"
	"strings"
)

const (
	OnDuKU = "ON DUPLICATE KEY UPDATE "
)

type Flags int

const (
	_ Flags = iota
	Insert
	Update
	OnDu
)

var FlagsName = []string{
	"insert",
	"update",
	"onDu",
}

func(f Flags)String()string{
	return FlagsName[f]
}

func(f Flags)Int()int{
	return int(f)
}

//CoverReqInfo request info
type CoverReqInfo struct {
	Type Flags
	Table string
	StructInfo interface{}
}

//CoverRespInfo resp info
type CoverRespInfo struct {
	Insert struct{
		Sql string
		Params []interface{}
	}
	Update struct{
		Sql string
		Params []interface{}
	}
	OnDu struct{
		Sql string
		Params []interface{}
	}
}

//GetSQL Get the insert, update, or MySQL on duplicate key update statements of the structure
func GetSQL(coverIn CoverReqInfo)(resp *CoverRespInfo,err error){
	t := reflect.ValueOf(coverIn.StructInfo)
	if err = checkParams(coverIn.Table,t);err != nil{
		return
	}
	if t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct{
		return nil,errors.New("params is not struct")
	}
	ty := t.Type()
	fieldNum := t.NumField()
	resp = &CoverRespInfo{
		Insert: struct {
			Sql    string
			Params []interface{}
		}{},
		Update: struct {
			Sql    string
			Params []interface{}
		}{},
		OnDu: struct {
			Sql    string
			Params []interface{}
		}{},
	}
	resp.Insert.Sql = "INSERT INTO " + coverIn.Table + "("
	resp.Update.Sql = "UPDATE SET " + coverIn.Table
	resp.OnDu.Sql +=  coverIn.Table

	for i := 1; i < fieldNum; i++ {
		resp.Insert.Sql += "`" + strings.ToLower(ty.Field(i).Name)+"`,"
		resp.Update.Sql += " `" + strings.ToLower(ty.Field(i).Name)+"` = ?,"
		resp.OnDu.Sql += " `" + strings.ToLower(ty.Field(i).Name)+"` = ?,"
		resp.Insert.Params = append(resp.Insert.Params,pkg.StringFromAssertionFloat(t.Field(i).Interface()))
		resp.Update.Params = append(resp.Update.Params,pkg.StringFromAssertionFloat(t.Field(i).Interface()))
	}
	resp.Insert.Sql = strings.TrimRight(resp.Insert.Sql,",")
	resp.Update.Sql = strings.TrimRight(resp.Update.Sql,",")
	resp.OnDu.Sql = strings.TrimRight(resp.OnDu.Sql,",")
	resp.OnDu.Sql = resp.Insert.Sql + OnDuKU + strings.TrimRight(resp.OnDu.Sql,",")
	resp.OnDu.Params = append(resp.OnDu.Params,resp.Insert.Params...)
	resp.OnDu.Params = append(resp.OnDu.Params,resp.Update.Params...)
	resp.Insert.Sql += ") VALUES (?"
	for i := 1; i < fieldNum; i++ {
		resp.Insert.Sql += ",?"
	}
	resp.Insert.Sql += ")"
	return SwitchResp(coverIn.Type,resp) ,nil
}


//SwitchResp response
func SwitchResp(typ Flags,resp *CoverRespInfo)*CoverRespInfo{
	switch typ {
	case Insert:
		return &CoverRespInfo{Insert: resp.Insert}
	case Update:
		return &CoverRespInfo{Update: resp.Update}
	case OnDu:
		return &CoverRespInfo{OnDu: resp.Update}
	default:
		return resp
	}
}


//checkParams check args
func checkParams(table string,t reflect.Value)error{
	if table == ""{
		return errors.New("table is empty")
	}
	if t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct{
		return errors.New("params is not struct")
	}
	return nil
}