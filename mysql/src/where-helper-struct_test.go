package mysql

import (
	"testing"
	"time"
)

func BenchmarkGenWhereByStruct(b *testing.B) {
	type Tmp struct{
		Addr string `column:"and,addr,like*"`
		Desc string `column:"and,desc,like"`
		Job string`column:"and,job,*like"`
		Name string `column:"and,name,="`
		Sal float32 `column:"and,sal,>"`
		AgeMin int`column:"or,age,between"`
		AgeMax int `column:"or,age,between"`
		Start time.Time `column:"and,created,between"`
		Stop time.Time `column:"and,created,between"`
		Jump string `column:"-"`
		Ptr *string `column:"and,ptr,="`
		Ptr1 string `column:"and,pt1r,="`
		Num int64 `column:"or,num,="`
		Num1 *int64 `column:"or,num1,="`
		Float *float32 `column:"or,flo,="`
		Float1 float32 `column:"or,flo1,="`
	}

	te := ""
	var te1 int64 = 0
	var te2 float32 = 0
	var tmp = Tmp{
		Addr:"earth",
		Name:"ft",
		Sal:333,
		AgeMin:9,
		AgeMax:18,
		Desc:"happ",
		Job:"engineer",
		Jump:"jump",
		Ptr: &te,
		Ptr1: " ",
		Num1: &te1,
		Float: &te2,
		Float1: 89,

	}

	for i := 0; i < b.N; i++ {
		_,_ = GenWhereByStruct(tmp,"column")
	}


}
