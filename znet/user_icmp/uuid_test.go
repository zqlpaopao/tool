package user_icmp

import (
	"testing"
)

func BenchmarkInitPool(b *testing.B) {
	//InitPool(16, 16, 36, 36)
	//for i := 0; i < b.N; i++ {
	//	u := NewUUid()
	//	s := u.String()
	//	UUIDPool.Uuid.Put(u)
	//	b := src.String2Bytes(s)
	//	UUIDPool.Byte.Put(b)
	//}

}

/*
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/znet/user_icmp
BenchmarkInitPool-10             3160761               377.8 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/zqlpaopao/tool/znet/user_icmp     2.008s
*/
