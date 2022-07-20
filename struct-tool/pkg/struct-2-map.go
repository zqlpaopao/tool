package pkg

import (
	"errors"
	"reflect"
	"strings"
)

const Omitempty = "omitempty"

// StructSlice2Map
//s The structure to be converted can be a pointer
//tagName must When formulating the tag of structure conversion, there must be
//fields Specifies the field of structure conversion. If there is no field, it is full conversion
//s []*struct *[]struct []struct{*int...}
func StructSlice2Map(s interface{}, tagName string, fields ...string) ([]map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("struct cannot be empty")
	}
	if tagName == "" {
		return nil, errors.New("tagName parameter cannot be empty")
	}

	typeSliOf := reflect.TypeOf(s)
	valueSliOf := reflect.ValueOf(s)
	if typeSliOf.Kind() == reflect.Ptr {
		valueSliOf = valueSliOf.Elem()
	}
	if typeSliOf.Kind() != reflect.Slice {
		return nil, errors.New("slice is not being converted")
	}

	fieldsLen := len(fields)
	sLen := reflect.ValueOf(s).Len()
	if sLen == 0 {
		return nil, nil
	}

	result := make([]map[string]interface{}, sLen)
	for i := 0; i < sLen; i++ {
		elem := make(map[string]interface{}, fieldsLen)
		allElem := make(map[string]interface{}, sLen)
		for k := range fields {
			elem[fields[k]] = 0
		}
		valueOf := reflect.ValueOf(s).Index(i)
		if valueOf.Type().Kind() == reflect.Ptr {
			valueOf = reflect.ValueOf(s).Index(i).Elem()
		}
		typeOf := valueOf.Type()
		numField := valueOf.NumField()
		for i := 0; i < numField; i++ {
			keySl := typeOf.Field(i).Tag.Get(tagName)
			kSL := strings.Split(keySl,",")
			if len(kSL) < 1{
				continue
			}
			key := kSL[0]
			omitempty := ""
			if len(kSL) >= 2{
				omitempty = kSL[1]
			}
			value := valueOf.Field(i).Interface()
			if omitempty == Omitempty && (value == 0 || value =="" ||value == nil){
				continue
			}
			if valueOf.Field(i).Kind() == reflect.Ptr {
				if valueOf.Field(i).Elem() == (reflect.Value{}){
					continue
				}
				value = valueOf.Field(i).Elem().Interface()
			}
			if _, ok := elem[key]; !ok {
				allElem[key] = value
				continue
			}
			elem[key] = value
		}

		if fieldsLen > 0 && len(elem) > 0 {
			result[i] = elem
		} else if fieldsLen < 1 {
			result[i] = allElem
		}

	}
	if len(result) == 0 {
		result = nil
	}
	return result, nil
}



//Struct2Map Convert structure to map
//s The structure to be converted can be a pointer
//tagName must When formulating the tag of structure conversion, there must be
//fields Specifies the field of structure conversion. If there is no field, it is full conversion
func Struct2Map(s interface{}, tagName string, fields ...string) (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("struct cannot be empty")
	}
	if tagName == "" {
		return nil, errors.New("tagName parameter cannot be empty")
	}
	typeOf := reflect.TypeOf(s)
	valueOf := reflect.ValueOf(s)
	if typeOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("struct is not being converted")
	}
	fieldsLen := len(fields)
	numField := valueOf.NumField()

	allResult := make(map[string]interface{}, numField)
	result := make(map[string]interface{}, fieldsLen)
	if numField == 0 {
		return nil, errors.New("no fields in the struct")
	}
	for i := 0; i < numField; i++ {
		keySl := typeOf.Field(i).Tag.Get(tagName)
		kSL := strings.Split(keySl,",")
		if len(kSL) < 1{
			continue
		}
		key := kSL[0]
		omitempty := ""
		if len(kSL) >= 2{
			omitempty = kSL[1]
		}
		value := valueOf.Field(i).Interface()
		if omitempty == Omitempty && (value == 0 || value =="" ||value == nil){
			continue
		}
		if valueOf.Field(i).Kind() == reflect.Ptr {
			if valueOf.Field(i).Elem() == (reflect.Value{}){
				continue
			}
			value = valueOf.Field(i).Elem().Interface()
		}
		result[key] = 0
		if _, ok := result[key]; !ok {
			allResult[key] = value
			continue
		}
		result[key] = value
	}
	if len(result) == 0 && len(allResult) == 0 {
		return nil, nil
	}
	if fieldsLen > 0 {
		return result, nil
	}
	return allResult, nil
}

