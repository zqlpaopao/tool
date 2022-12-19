package pkg

import (
	"errors"
	"reflect"
	"unsafe"
)

type TypeModel uint

const (
	StringType TypeModel = iota

	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
	UintType

	Int8Type
	Int16Type
	Int32Type
	Int64Type
	IntType

	Float32Type
	Float64Type

	SliceStringType
	SliceIntType

	SliceUint8Type
	SliceUint16Type
	SliceUint32Type
	SliceUint64Type
	SliceUintType

	SliceInt8Type
	SliceInt16Type
	SliceInt32Type
	SliceInt64Type

	SliceFloat32Type
	SliceFloat64Type
)

var TypeValue = []string{
	"string",

	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uint",

	"int8",
	"int16",
	"int32",
	"int64",
	"int",

	"float32",
	"float64",

	"sliceString",
	"sliceInt",

	"sliceUint8",
	"sliceUint16",
	"sliceUint32",
	"sliceUint64",
	"sliceUint",

	"sliceInt8",
	"sliceInt16",
	"sliceInt32",
	"sliceInt64",
	"sliceInt",

	"sliceFloat32",
	"sliceFloat64",
}

func (t TypeModel) String() string {
	return TypeValue[t]
}

type (
	TagName struct {
		StructName string
		Type       TypeModel
		MapKey     string
		Offset     uintPtrDescriptor
		Fun        Func
	}
	modelFace struct {
		typ   unsafe.Pointer
		value unsafe.Pointer
	}
	uintPtrDescriptor uintptr
	Func              func(structPoint unsafe.Pointer, ti uintPtrDescriptor, val interface{})
)

var TypeMap = map[TypeModel]Func{
	StringType: makeString,

	SliceStringType: makeSliceString,

	Uint8Type:  makeUint8,
	Uint16Type: makeUint16,
	Uint32Type: makeUint32,
	Uint64Type: makeUint64,
	UintType:   makeUint,

	Int8Type:  makeInt8,
	Int16Type: makeInt16,
	Int32Type: makeInt32,
	Int64Type: makeInt64,
	IntType:   makeInt,

	Float32Type: makeFloat32,
	Float64Type: makeFloat64,

	SliceUint8Type:  makeSliceUint8,
	SliceUint16Type: makeSliceUint16,
	SliceUint32Type: makeSliceUint32,
	SliceUint64Type: makeSliceUint64,
	SliceUintType:   makeSliceUint,

	SliceInt8Type:  makeSliceInt8,
	SliceInt16Type: makeSliceInt16,
	SliceInt32Type: makeSliceInt32,
	SliceInt64Type: makeSliceInt64,
	SliceIntType:   makeSliceInt,

	SliceFloat32Type: makeSliceFloat32,
	SliceFloat64Type: makeSliceFloat64,
}

func DescribeStructUnsafePointer(in interface{}, tagName []*TagName) (err error) {
	typ := reflect.TypeOf(in)
	tagNameRef := make([]*TagName, 0, len(tagName))
	if typ.Kind() != reflect.Ptr {
		err = errors.New("you must pass in a pointer")
		return
	}
	if typ.Elem().Kind() != reflect.Struct {
		err = errors.New("you must pass in a pointer to a struct")
		return
	}
	for k, v := range tagName {
		f, ok := typ.Elem().FieldByName(v.StructName)
		if !ok {
			continue
		}
		tagName[k].Offset = uintPtrDescriptor(f.Offset)
		fun, ok := TypeMap[tagName[k].Type]
		if !ok {
			continue
		}
		tagName[k].Fun = fun
		tagNameRef = append(tagNameRef, v)
	}
	tagName = nil
	tagName = tagNameRef
	return
}

func Map2StructOver(in interface{}, tagName []*TagName, valMap map[string]interface{}) (err error) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	for i := 0; i < len(tagName); i++ {
		tagName[i].Fun(structPtr, tagName[i].Offset, valMap[tagName[i].MapKey])
	}
	return
}

//******************************************* string **********************************//
func makeString(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*string)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(string)
}

//******************************************* uint8-64 **********************************//
func makeUint8(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*uint8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(uint8)
}

func makeUint16(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*uint16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(uint16)
}

func makeUint32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*uint32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(uint32)
}

func makeUint64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*uint64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(uint64)
}

func makeUint(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*uint)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(uint)
}

//******************************************* int8-64 **********************************//
func makeInt8(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*int8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(int8)
}

func makeInt16(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*int16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(int16)
}

func makeInt32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*int32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(int32)
}

func makeInt64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*int64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(int64)
}

func makeInt(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*int)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(int)
}

//******************************************* float **********************************//
func makeFloat32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(float32)
}
func makeFloat64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.(float64)
}

//******************************************* slice string **********************************//

func makeSliceString(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]string)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]string)
}

//******************************************* slice int **********************************//
func makeSliceInt(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]int)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]int)
}

//*******************************************  slice uint8-64 **********************************//
func makeSliceUint8(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]uint8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]uint8)
}

func makeSliceUint16(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]uint16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]uint16)
}

func makeSliceUint32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]uint32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]uint32)
}

func makeSliceUint64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]uint64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]uint64)
}

func makeSliceUint(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]uint)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]uint)
}

//*******************************************  slice int8-64 **********************************//
func makeSliceInt8(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]int8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]int8)
}

func makeSliceInt16(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]int16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]int16)
}

func makeSliceInt32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]int32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]int32)
}

func makeSliceInt64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]int64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]int64)
}

//*******************************************  slice float **********************************//
func makeSliceFloat32(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]float32)
}
func makeSliceFloat64(structPtr unsafe.Pointer, ti uintPtrDescriptor, val interface{}) {
	*(*[]float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val.([]float64)
}
