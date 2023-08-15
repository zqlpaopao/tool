package ping

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

type Pools struct {
	b chan []byte
}

func (p *Pools) GetByte() (res []byte) {
	select {
	case res, _ = <-p.b:
		return
	default:
		return
	}
}

func (p *Pools) SetByte(res []byte) {
	select {
	case p.b <- res:
		return
	default:
		return
	}
}

func BenchmarkGetByteByPoolGet(b *testing.B) {
	var p = &Pools{b: make(chan []byte, 100)}

	for i := 0; i < b.N; i++ {
		_ = p.GetByte()

		//p.SetByte(res)
	}
}

func BenchmarkGetByteByPoolSet(b *testing.B) {
	var p = &Pools{b: make(chan []byte, 100)}

	for i := 0; i < b.N; i++ {
		res := []byte{}

		p.SetByte(res)
	}
}

type PoolsSync struct {
	b *sync.Pool
}

func BenchmarkPoolsSyncGet(b *testing.B) {
	var p = &PoolsSync{b: &sync.Pool{New: func() any {
		return []byte{}
	}}}

	for i := 0; i < b.N; i++ {
		res := p.b.Get()
		if _, ok := res.([]byte); !ok {
			continue
		}
		//p.b.Put(res)
	}
}

func BenchmarkPoolsSyncSet(b *testing.B) {
	var p = &PoolsSync{b: &sync.Pool{New: func() any {
		return []byte{}
	}}}

	for i := 0; i < b.N; i++ {
		res := []byte{}
		//if _, ok := res.([]byte); !ok {
		//	continue
		//}
		p.b.Put(res)
	}
}

type PoolsZeroMap struct {
	b *zeropool.Pool[[]byte]
}

func BenchmarkPoolsZeroMapGet(b *testing.B) {
	var p = zeropool.New[[]byte](func() []byte {
		return []byte{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()

		//p.Put(res)
	}
}

func BenchmarkPoolsZeroMapSet(b *testing.B) {
	var p = zeropool.New[[]byte](func() []byte {
		return []byte{}
	})

	for i := 0; i < b.N; i++ {
		res := []byte{}

		p.Put(res)
	}
}

//
///*
//
//	go test -bench=. -benchmem
//goos: darwin
//goarch: arm64
//pkg: github.com/zqlpaopao/tool/batch_ping_agent/pkg
//BenchmarkGetByteByPoolGet-10            459246643                2.567 ns/op           0 B/op          0 allocs/op
//BenchmarkGetByteByPoolSet-10            515846056                2.301 ns/op           0 B/op          0 allocs/op
//BenchmarkPoolsSyncGet-10                25145130                48.55 ns/op           24 B/op          1 allocs/op
//BenchmarkPoolsSyncSet-10                33763821                32.06 ns/op           53 B/op          1 allocs/op
//BenchmarkPoolsZeroMapGet-10             18812818                61.75 ns/op           51 B/op          1 allocs/op
//BenchmarkPoolsZeroMapSet-10             19779064                59.22 ns/op           52 B/op          1 allocs/op
//PASS
//ok      github.com/zqlpaopao/tool/batch_ping_agent/pkg  7.911s
//
//
//*/
