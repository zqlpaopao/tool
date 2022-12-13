package pkg

import (
	"fmt"
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

//自定义选择版本
type tagName struct {
	Type string
	key  string
}

type uintPtrDescriptor uintptr

func Map2StructOver(m map[string]interface{}, in interface{}, tagName map[string]*tagName) (err error) {
	typ := reflect.TypeOf(in)
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("you must pass in a pointer")
	}
	if typ.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("you must pass in a pointer to a struct")
	}

	for k, v := range tagName {
		f, ok := typ.Elem().FieldByName(k)
		if !ok {
			continue
		}
		if _, ok := m[v.key]; !ok {
			continue
		}
		SwitchType(uintPtrDescriptor(f.Offset), in, v, m[v.key])
	}
	return
}

func SwitchType(offset uintPtrDescriptor, in interface{}, tag *tagName, val interface{}) {
	if tag.Type == StringType.String() {
		makeString(in, offset, val.(string))
		//
	} else if tag.Type == IntType.String() {
		makeInt(in, offset, val.(int))
		//
	} else if tag.Type == SliceStringType.String() {
		makeSliceString(in, offset, val.([]string))
	} else if tag.Type == SliceIntType.String() {
		makeSliceInt(in, offset, val.([]int))

	} else if tag.Type == Int64Type.String() {
		makeInt64(in, offset, val.(int64))
	} else if tag.Type == Uint8Type.String() {
		makeUint8(in, offset, val.(uint8))
	} else if tag.Type == Uint16Type.String() {
		makeUint16(in, offset, val.(uint16))
	} else if tag.Type == Uint32Type.String() {
		makeUint32(in, offset, val.(uint32))
	} else if tag.Type == Uint64Type.String() {
		makeUint64(in, offset, val.(uint64))
	} else if tag.Type == UintType.String() {
		makeUint(in, offset, val.(uint))
		//
	} else if tag.Type == Int8Type.String() {
		makeInt8(in, offset, val.(int8))
	} else if tag.Type == Int16Type.String() {
		makeInt16(in, offset, val.(int16))
	} else if tag.Type == Int32Type.String() {
		makeInt32(in, offset, val.(int32))
	} else if tag.Type == Float32Type.String() {
		makeFloat32(in, offset, val.(float32))
	} else if tag.Type == Float64Type.String() {
		makeFloat64(in, offset, val.(float64))
	} else if tag.Type == SliceUint8Type.String() {
		makeSliceUint8(in, offset, val.([]uint8))
	} else if tag.Type == SliceUint16Type.String() {
		makeSliceUint16(in, offset, val.([]uint16))
	} else if tag.Type == SliceUint32Type.String() {
		makeSliceUint32(in, offset, val.([]uint32))
	} else if tag.Type == SliceUint64Type.String() {
		makeSliceUint64(in, offset, val.([]uint64))
	} else if tag.Type == SliceUintType.String() {
		makeSliceUint(in, offset, val.([]uint))
		//
	} else if tag.Type == SliceInt8Type.String() {
		makeSliceInt8(in, offset, val.([]int8))
	} else if tag.Type == SliceInt16Type.String() {
		makeSliceInt16(in, offset, val.([]int16))
	} else if tag.Type == SliceInt32Type.String() {
		makeSliceInt32(in, offset, val.([]int32))
	} else if tag.Type == SliceInt64Type.String() {
		makeSliceInt64(in, offset, val.([]int64))
		//
	} else if tag.Type == SliceFloat32Type.String() {
		makeSliceFloat32(in, offset, val.([]float32))
	} else if tag.Type == SliceFloat64Type.String() {
		makeSliceFloat64(in, offset, val.([]float64))
		//
	}

}

//******************************************* string **********************************//
func makeString(in interface{}, ti uintPtrDescriptor, val string) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*string)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//******************************************* uint8-64 **********************************//
func makeUint8(in interface{}, ti uintPtrDescriptor, val uint8) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*uint8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeUint16(in interface{}, ti uintPtrDescriptor, val uint16) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*uint16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeUint32(in interface{}, ti uintPtrDescriptor, val uint32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*uint32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeUint64(in interface{}, ti uintPtrDescriptor, val uint64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*uint64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeUint(in interface{}, ti uintPtrDescriptor, val uint) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*uint)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//******************************************* int8-64 **********************************//
func makeInt8(in interface{}, ti uintPtrDescriptor, val int8) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*int8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeInt16(in interface{}, ti uintPtrDescriptor, val int16) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*int16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeInt32(in interface{}, ti uintPtrDescriptor, val int32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*int32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeInt64(in interface{}, ti uintPtrDescriptor, val int64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*int64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeInt(in interface{}, ti uintPtrDescriptor, val int) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*int)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//******************************************* float **********************************//
func makeFloat32(in interface{}, ti uintPtrDescriptor, val float32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}
func makeFloat64(in interface{}, ti uintPtrDescriptor, val float64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//******************************************* slice string **********************************//

func makeSliceString(in interface{}, ti uintPtrDescriptor, val []string) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]string)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//******************************************* slice int **********************************//
func makeSliceInt(in interface{}, ti uintPtrDescriptor, val []int) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]int)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//*******************************************  slice uint8-64 **********************************//
func makeSliceUint8(in interface{}, ti uintPtrDescriptor, val []uint8) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]uint8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceUint16(in interface{}, ti uintPtrDescriptor, val []uint16) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]uint16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceUint32(in interface{}, ti uintPtrDescriptor, val []uint32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]uint32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceUint64(in interface{}, ti uintPtrDescriptor, val []uint64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]uint64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceUint(in interface{}, ti uintPtrDescriptor, val []uint) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]uint)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//*******************************************  slice int8-64 **********************************//
func makeSliceInt8(in interface{}, ti uintPtrDescriptor, val []int8) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]int8)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceInt16(in interface{}, ti uintPtrDescriptor, val []int16) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]int16)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceInt32(in interface{}, ti uintPtrDescriptor, val []int32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]int32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

func makeSliceInt64(in interface{}, ti uintPtrDescriptor, val []int64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]int64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}

//*******************************************  slice float **********************************//
func makeSliceFloat32(in interface{}, ti uintPtrDescriptor, val []float32) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]float32)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}
func makeSliceFloat64(in interface{}, ti uintPtrDescriptor, val []float64) {
	structPtr := (*modelFace)(unsafe.Pointer(&in)).value
	*(*[]float64)(unsafe.Pointer(uintptr(structPtr) + uintptr(ti))) = val
}
