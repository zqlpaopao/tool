package mysql

import (
	"errors"
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"reflect"
	"strings"
)

const (
	OnDuKU = " ON DUPLICATE KEY UPDATE "
	Tag = "cv"
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
	var params []interface{}
	resp.Insert.Sql,  resp.Update.Sql,resp.OnDu.Sql,params = getSqlAndParams(coverIn.Table,t)
	resp.Insert.Params = params
	resp.Update.Params = params
	resp.OnDu.Params = params
	resp.OnDu.Params = append(resp.OnDu.Params,params...)
	return SwitchResp(coverIn.Type,resp) ,nil
}

//getSqlAndParams get sql & params
func getSqlAndParams(table string,t reflect.Value)(insert,update,onDu string,params []interface{}){
	ty := t.Type()
	var (
		insertSql,updateSql,onDuSql strings.Builder
	)
	insertSql.WriteString("INSERT INTO "+table + "(")
	updateSql.WriteString("UPDATE SET")
	onDuSql.WriteString(table)
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		if ty.Field(i).Tag.Get(Tag) == ""{
			insertSql.WriteString("`" + pkg.SnakeString(ty.Field(i).Name)+"`,")
			updateSql.WriteString(" `" + pkg.SnakeString(ty.Field(i).Name)+"` = ?,")
			onDuSql.WriteString(" `" + pkg.SnakeString(ty.Field(i).Name)+"` = ?,")
		}else {
			insertSql.WriteString("`" + ty.Field(i).Tag.Get("cv")+"`,")
			updateSql.WriteString(" `" + ty.Field(i).Tag.Get("cv")+"` = ?,")
			onDuSql.WriteString(" `" + ty.Field(i).Tag.Get("cv")+"` = ?,")
		}
		params = append(params,pkg.StringFromAssertionFloat(t.Field(i).Interface()))
	}
	insert,update,onDu = insertSql.String(),updateSql.String(),onDuSql.String()

	insert = strings.TrimRight(insert,",")
	update = strings.TrimRight(update,",")
	insert += ") VALUES (?"
	for i := 1; i < fieldNum; i++ {
		insert += ",?"
	}
	insert += ")"
	onDu = insert + OnDuKU + strings.TrimRight(onDu,",")

	return
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