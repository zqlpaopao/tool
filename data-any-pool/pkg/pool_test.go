package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

func BenchmarkPoolGet(b *testing.B) {
	var p = &Pool[[]byte]{
		Data: make(chan *[]byte, 50),
		New: func() *[]byte {
			return &[]byte{}
		},
	}

	for i := 0; i < b.N; i++ {
		_ = p.Get()

		//p.SetByte(res)
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = &Pool[[]byte]{
		Data: make(chan *[]byte, 50),
		New: func() *[]byte {
			return &[]byte{}
		},
	}

	for i := 0; i < b.N; i++ {
		p.Put(&[]byte{})

	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = &Pool[[]byte]{
		Data: make(chan *[]byte, 50),
		New: func() *[]byte {
			return &[]byte{}
		},
	}

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return []byte{}
	}}

	for i := 0; i < b.N; i++ {
		res := p.Get()
		if _, ok := res.([]byte); !ok {
			continue
		}
		//p.b.Put(res)
	}
}

func BenchmarkSyncPoolSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return []byte{}
	}}

	for i := 0; i < b.N; i++ {
		res := []byte{}
		p.Put(res)
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return []byte{}
	}}

	for i := 0; i < b.N; i++ {
		res := p.Get()
		if _, ok := res.([]byte); !ok {
			continue
		}
		p.Put(res)
	}
}

func BenchmarkPoolsZeroGet(b *testing.B) {
	var p = zeropool.New[[]byte](func() []byte {
		return []byte{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()

		//p.Put(res)
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[[]byte](func() []byte {
		return []byte{}
	})

	for i := 0; i < b.N; i++ {
		res := []byte{}

		p.Put(res)
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[[]byte](func() []byte {
		return []byte{}
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()

		p.Put(res)
	}
}

/*
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/xxx/tool/data-any-pool/pkg
BenchmarkPoolGet-10             53016801                22.35 ns/op           24 B/op          1 allocs/op
BenchmarkSyncPoolGet-10         24650505                48.78 ns/op           24 B/op          1 allocs/op
BenchmarkPoolsZeroGet-10        18879105                61.91 ns/op           51 B/op          1 allocs/op

BenchmarkPoolSet-10             56471364                21.35 ns/op           24 B/op          1 allocs/op
BenchmarkSyncPoolSet-10         36137818                32.39 ns/op           52 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10        19359301                61.38 ns/op           49 B/op          1 allocs/op

BenchmarkPoolGetSet-10          43514325                27.35 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10      135167949                8.479 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10     69344953                17.26 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/xxx/tool/data-any-pool/pkg     12.094s



*/
