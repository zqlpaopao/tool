package mysql

import (
	"fmt"
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//GenWhereByStruct Get the where SQL according to the structure
func GenWhereByStruct(in interface{},tagName string)(string,[]interface{}){
	vValue :=reflect.ValueOf(in)
	if !checkArgsStr(vValue){
		return "",nil
	}
	vType :=reflect.TypeOf(in)
	var tagTmp =""
	var whereMap = make([][]string,0)
	var args = make([]interface{},0)
	var b bool
	for i:=0;i<vValue.NumField();i++{
		tagTmp = vType.Field(i).Tag.Get(tagName)
		b,_ = pkg.StringFromAssertionFloatNotZero(vType.Field(i).Type.String(),vValue.Field(i))
		if tagTmp =="-"||tagTmp=="" || !b{
			continue
		}
		fmt.Println(vType.Field(i).Name)
		cons :=strings.Split(tagTmp,",")
		if !IfZero(vValue.Field(i).Interface()) {
			var al  = vValue.Field(i)
			if vValue.Field(i).Kind()== reflect.Ptr{
				al = vValue.Field(i).Elem()
			}
			if cons[2]=="*like"{
				cons[2] = "like"
				args = append(args, "%"+al.Interface().(string))
			}else if cons[2]=="like*"{
				cons[2] = "like"
				args = append(args, al.Interface().(string)+"%")
			}else if cons[2]=="*like*" || cons[2]=="like"{
				cons[2] = "like"
				args = append(args, "%"+al.Interface().(string)+"%")
			}else{
				args = append(args, al.Interface())
			}
			if len(whereMap)==0 {
				whereMap = append(whereMap,[]string{
					"",cons[1],cons[2],
				})
			}else{
				whereMap = append(whereMap,[]string{
					cons[0],cons[1],cons[2],
				})
			}
			if cons[2] == "between"{
				i++
				args = append(args,al.Interface())
			}
		}
	}
	where :=GenWhere(whereMap)
	return where,args
}

//checkArgsStr Check whether it is a structure
func checkArgsStr(t reflect.Value) bool{
	if t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct{
		return  false
	}
	return true
}

//GenWhere generate where through a where [][]string
func GenWhere(whereMap [][]string)string {
	rs:=""
	if len(whereMap) != 0 {
		rs = rs + " where "
		for _, v := range whereMap {
			//v[0]表示性质，and 还是or,v[1]表示field，比如name，age,v[2]表示条件符号,=,>,<,<>,like
			if v[2] == "between" {
				rs = rs + " " + v[0] + " " + v[1] + " " + "between" + " " + "?" + " " + "and" + " " + "?" + " "
				continue
			}
			if v[2] == "in" {
				rs = rs + " " + v[0] + " " + v[1] + " " + "in" + " " +v[3]
				continue
			}
			rs = rs + " " + v[0] + " " + v[1] + " " + v[2] + " " + "?"
		}
	}
	return rs
}



//IfZero Judge whether it is 0
func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, float64, int32, int16, int64, float32:
		if v == 0 {
			return true
		}
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
	return false
}

func RemoveZero(slice []interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if IfZero(v) {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveZero(slice)
			break
		}
	}
	return slice
}
//ReplaceQuestionToDollar 将sql语句中的?转换成$i
func ReplaceQuestionToDollar(sql string) string {
	var temp = 1
	start := 0
	var i = 0
L:
	for i = start; i < len(sql); i++ {
		if string(sql[i]) == "?" {
			sql = string(sql[:i]) + "$" + strconv.Itoa(temp) + string(sql[i+1:])
			temp++
			start = i + 2
			goto L
		}

		if i == len(sql)-1 {
			return sql
		}
	}
	return sql
}
