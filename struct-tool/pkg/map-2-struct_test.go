package pkg

import (
	//"fmt"
	//jsoniter "github.com/json-iterator/go"
	//"strconv"
	//"test/reflect-cache/pkg"
	"testing"
)

func JsonMa(p interface{}) {
	//var res map[string]interface{}
	//k,_:=jsoniter.Marshal(p)
	//err := jsoniter.Unmarshal(k, &res)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

type P struct {
	Name string `json:"name"`
	Id float64 `json:"id"`
	Sex string `json:"sex"`
}


func TestSliceMapToSliceStruct(t *testing.T) {
	a := []map[string]interface{}{
		{"id": 213.3, "name": "zhangsan", "sex": "男"},
	}
	var s []*P

	err := SliceMapToSliceStruct(a, &s, "json")
	if err != nil {
		t.Log(err)
	}
	t.Log()
}


func BenchmarkStructToMap(b *testing.B)  {
	b.ResetTimer()
	a := []map[string]interface{}{
		{"Id": 213, "Name": "zhaoliu", "Sex": "男"},
		{"Id": 56, "Name": "zhangsan", "Sex": "男"},
		{"Id": 7, "Name": "lisi", "Sex": "女"},
		{"Id": 978, "Name": "wangwu", "Sex": "男"},
	}
	var s []*P
	for i:=0; i<b.N; i++{
		//_, _ = StructToMap(p)
		//JsonMa(p)
		SliceMapToSliceStruct(a, &s, "json")
	}
}
func BenchmarkStructToMap1(b *testing.B)  {
	b.ResetTimer()
	a := map[string]interface{}{
		"Id": 213, "Name": "zhaoliu", "Sex": "男",
	}
	//{"Id": 213, "Name": "zhaoliu", "Sex": "男"},
	//{"Id": 56, "Name": "zhangsan", "Sex": "男"},
	//{"Id": 7, "Name": "lisi", "Sex": "女"},
	//{"Id": 978, "Name": "wangwu", "Sex": "男"},

	var s P
	for i:=0; i<b.N; i++{
		//_, _ = StructToMap(p)
		//JsonMa(p)
		if err := MapToStruct(a, &s, "json");err != nil{
			b.Fatal(err)
		}
	}
}


func BenchmarkStructToMap12(b *testing.B) {
	a := map[string]interface{}{
		"id": 213, "name": "zhaoliu", "sex": "男",
	}
	//
	//m = m

	//Name string `json:"name"`
	//Id float64 `json:"id"`
	//Sex string `json:"sex"`

	var tagNames = []*TagName{
		{
			StructName: "Name",
			Type: "string",
			MapKey:  "name",
		},
		{
			StructName: "Id",
			Type: "int",
			MapKey:  "id",
		},
		{
			StructName: "Sex",
			Type: "string",
			MapKey:  "sex",
		},
	}

	b.StartTimer()
	err := DescribeStructUnsafePointer((*P)(nil),tagNames)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		//a["str"] = strconv.Itoa(i+100)
		var structs P
		//pkg.MapToStruct()
		if err := Map2StructOver(&structs,tagNames,a);nil != err{
			b.Fatal(err)
		}
		//fmt.Println(structs)
		//fmt.Println(structs.StrPtr)
		//fmt.Println(structs.StrPt)
	}
}

/*
go test -bench=^BenchmarkStructToMap1  -v --benchmem
=== RUN   TestSliceMapToSliceStruct
    map-2-struct_test.go:37:
--- PASS: TestSliceMapToSliceStruct (0.00s)
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/struct-tool/pkg
BenchmarkStructToMap1
BenchmarkStructToMap1-10         3271748               350.3 ns/op            72 B/op          6 allocs/op
BenchmarkStructToMap12
BenchmarkStructToMap12-10       27577093                43.74 ns/op           48 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/struct-tool/pkg       3.184s



*/
