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

func map2Struct() {
	type User struct {
		Name string
		ID   int
	}

	m := map[string]interface{}{}
	m["Name"] = "test"
	m["ID"] = 100
	user := new(User)
	err := pkg.Map2Struct(user, m)
	fmt.Println(err)
	fmt.Printf("%#v", user)
}
