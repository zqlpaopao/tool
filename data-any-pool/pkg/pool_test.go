package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 100

type Model *Model1
type Model1 struct {
	Name    string
	Age     string
	Sex     int
	Address map[string]string
}

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		100,
		func(size int) Model {
			return &Model1{}
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return &Model1{}
		})

	for i := 0; i < b.N; i++ {
		p.Put(&Model1{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return &Model1{}
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return &Model1{}
	}}

	for i := 0; i < b.N; i++ {
		res := p.Get()
		if _, ok := res.(Model); !ok {
			continue
		}
		//p.b.Put(res)
	}
}

func BenchmarkSyncPoolSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return &Model1{}
	}}
	for i := 0; i < b.N; i++ {
		p.Put(&Model1{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return &Model1{}
	}}

	for i := 0; i < b.N; i++ {
		res := p.Get()
		if _, ok := res.(Model); !ok {
			continue
		}
		p.Put(res)
	}
}

func BenchmarkPoolsZeroGet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		p.Put(&Model1{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		p.Put(&Model1{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return &Model1{}
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}
