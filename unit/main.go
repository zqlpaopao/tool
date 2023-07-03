package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/unit/pkg"
)

func main() {
	//fmt.Println(pkg.IBytes(82814982, 3))
	//fmt.Println(pkg.Times(82854982, 3))
	////fmt.Println(pkg.Times(1750, 2))
	fmt.Println(1755, pkg.TimesHalfAdjust(1755, 2))
	fmt.Println(1759, pkg.TimesHalfAdjust(1759, 2))
	fmt.Println(1756, pkg.TimesHalfAdjust(1756, 2))
	fmt.Println(1750, pkg.TimesHalfAdjust(1750, 2))
	fmt.Println(1751, pkg.TimesHalfAdjust(1751, 2))
	fmt.Println(17517, pkg.TimesHalfAdjust(1751789, 3))
	//fmt.Println(pkg.TimesHalfAdjust(1755, 2))
	//
	//fmt.Println(pkg.IBytes(82854982, 3))
	//fmt.Println(pkg.IBytes(82854982, 2))
	//fmt.Println(pkg.IBytes(82854982, 0))
	//
	//fmt.Println(pkg.Bytes(82854982, 3))
	//fmt.Println(pkg.Bytes(82854982, 2))
	//fmt.Println(pkg.Bytes(82854982, 0))

	//Bytes(82854982) -> 83 MB
}
