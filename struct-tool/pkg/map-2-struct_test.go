package pkg

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func JsonMa(p interface{}) {
	var res map[string]interface{}
	k,_:=jsoniter.Marshal(p)
	err := jsoniter.Unmarshal(k, &res)
	if err != nil {
		fmt.Println(err)
	}
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
