package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/struct-tool/pkg"
)

type Mes struct {
	Name string `cover:"name"`
	Age  *int   `cover:"age"`
	Sex  int    `cover:"sex"`
}

func main() {
	map2Struct()
}

func struct2map() {
	age := 5
	var mes = Mes{
		Name: "Name",
		Age:  &age,
		Sex:  1,
	}

	var mes1 []*Mes
	mes1 = append(mes1, &mes)
	//fmt.Println(Struct2Map(mes))
	//fmt.Println(Struct2Map(mes,"cover",[]string{"name","age"}...))
	fmt.Println(pkg.StructSlice2Map(mes1, "cover", []string{}...))
}


type P struct {
	Name string `cover:"name"`
	Id float64 `json:"id"`
	Sex string `json:"sex"`
}

func map2Struct() {
	type User struct {
		Name string
		ID   int
	}

	a := []map[string]interface{}{
		{"id": 213.3, "name": "zhaoliu", "sex": "男"},
	}
	var s []*P

	err := pkg.SliceMapToSliceStruct(a, &s, "cover")
	if err != nil {
		panic(err)
	}
	for _ ,v := range s{
		fmt.Printf("%#v",v)

	}
}
