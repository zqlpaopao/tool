package main

import (
	"fmt"
	"os"
	"testing"
)
func TestMain(m *testing.M) {
	fmt.Println("write setup code here...") // 测试之前的做一些设置
	// 如果 TestMain 使用了 flags，这里应该加上flag.Parse()
	retCode := m.Run()                         // 执行测试
	fmt.Println("write teardown code here...") // 测试之后做一些拆卸工作
	os.Exit(retCode)                           // 退出测试
}

func TestLog(t *testing.T) {
	Log()
}
func TestAsyncLog(t *testing.T) {
	AsyncLog()

}

func BenchmarkLog(b *testing.B) {
	for i := 0; i < 10000; i++ {
		Log()
	}
}

func BenchmarkAsyncLog(b *testing.B) {
	for i := 0; i < 10000; i++ {
		AsyncLog()
	}
}
