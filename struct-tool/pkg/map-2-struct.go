package pkg

import (
	"errors"
	"reflect"
)

//判断是否是int大类
func isIntType(k reflect.Kind, values interface{}) (bool, interface{}) {
	switch k {
	case reflect.Uint8:
		return true, uint8(values.(int64))
	case reflect.Uint16:
		return true, uint16(values.(int64))
	case reflect.Uint32:
		return true, uint32(values.(int64))
	case reflect.Uint64:
		return true, uint64(values.(int64))
	case reflect.Int8:
		return true, int8(values.(int64))
	case reflect.Int16:
		return true, int16(values.(int64))
	case reflect.Int32:
		return true, int32(values.(int64))
	case reflect.Int64:
		return true, values
	default:
		return false, nil
	}
}

//判断是否float大类
func isFloatType(k reflect.Kind, values interface{}) (bool, interface{}) {
	switch k {
	case reflect.Float32:
		return true, float32(values.(float64))
	case reflect.Float64:
		return true, values
	default:
		return false, nil
	}
}

//SliceMapToSliceStruct slice的map转slice的struct
func SliceMapToSliceStruct(mList []map[string]interface{}, model interface{}, tagName string) (err error) {
	val := reflect.Indirect(reflect.ValueOf(model))
	typ := val.Type()
	items := make([]reflect.Value, 0, len(mList))
	for _, r := range mList {
		mVal := reflect.Indirect(reflect.New(typ.Elem().Elem())).Addr()
		err = MapToStruct(&r, mVal.Interface(), tagName)
		if err != nil {
			return err
		}
		items = append(items, mVal)
	}
	values := reflect.Append(val, items...)
	val.Set(values)
	return err
}

//MapToStruct map转struct，map的key和struct的key有小写要求
func MapToStruct(m, s interface{}, tagName string) error {
	//如果传递的是指针，直接解引用
	mVal := reflect.Indirect(reflect.ValueOf(m))
	mValType := mVal.Type()
	if mValType.Kind() != reflect.Map {
		return errors.New("类型错误")
	}

	sVal := reflect.Indirect(reflect.ValueOf(s))
	sValType := sVal.Type()
	if sValType.Kind() != reflect.Struct {
		return errors.New("类型错误")
	}
	for i := 0; i < sVal.NumField(); i++ {
		//需要先判断是否有tag标签
		sField := sValType.Field(i)
		var name string
		if tagName != "" {
			name = sField.Tag.Get(tagName)
		}
		if name == "" {
			name = sField.Name
		}
		mKey := mVal.MapIndex(reflect.ValueOf(name))
		if !mKey.IsValid() {
			continue
		}
		if mKey.IsZero() {
			continue
		}
		//由于从map中获取的int值都是默认int，float默认float64，因此需要做特殊处理
		values := mKey.Elem().Interface()
		//fmt.Println(values, reflect.TypeOf(values))
		ok := false
		switch reflect.TypeOf(values).Kind() {
		case reflect.Int:
			ok, values = isIntType(sField.Type.Kind(), int64(values.(int)))
			if !ok {
				continue
			}
		case reflect.Int32:
			ok, values = isIntType(sField.Type.Kind(), int64(values.(int32)))
			if !ok {
				continue
			}
		case reflect.Int64:
			ok, values = isIntType(sField.Type.Kind(), values)
			if !ok {
				continue
			}
		case reflect.Float64:
			ok, values = isFloatType(sField.Type.Kind(), values)
			if !ok {
				continue
			}
		default:
			if mKey.Elem().Type() != sField.Type {
				continue
			}
		}
		sValField := sVal.Field(i)
		if sValField.CanSet() {
			sValField.Set(reflect.ValueOf(values))
		}
	}
	return nil
}

//StructToMap map转struct，可以指定map的value类型和struct的key有小写要求
func StructToMap(s interface{}) (error, *map[string]interface{}) {
	sVal := reflect.Indirect(reflect.ValueOf(s))
	sType := sVal.Type()
	if sType.Kind() != reflect.Struct {
		return errors.New("类型错误"), nil
	}
	//读取map的key和value的类型，类型不同，不能转换
	res := make(map[string]interface{}, sVal.NumField())
	for i := 0; i < sVal.NumField(); i++ {
		key := sType.Field(i).Name
		val := sVal.FieldByIndex([]int{i})
		if !val.IsValid() {
			continue
		}
		res[key] = val.Interface()
	}
	return nil, &res
}
