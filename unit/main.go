package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/byte_unit/pkg"
)

func main() {
	//fmt.Println(pkg.IBytes(82814982, 3))
	fmt.Println(pkg.Times(82854982, 3))
	fmt.Println(pkg.Times(82814982, 2))

	fmt.Println(pkg.IBytes(82854982, 3))
	fmt.Println(pkg.IBytes(82854982, 2))
	fmt.Println(pkg.IBytes(82854982, 0))

	fmt.Println(pkg.Bytes(82854982, 3))
	fmt.Println(pkg.Bytes(82854982, 2))
	fmt.Println(pkg.Bytes(82854982, 0))

	//Bytes(82854982) -> 83 MB
}
