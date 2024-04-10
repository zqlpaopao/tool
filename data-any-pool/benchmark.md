# []byte-3

```
package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 3

type Model []byte

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	}}
	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

```



```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     79894804                14.41 ns/op            3 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 20696583                57.85 ns/op           27 B/op          2 allocs/op
BenchmarkPoolsZeroGet-10                16648543                73.38 ns/op           57 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              20583882                58.17 ns/op           27 B/op          2 allocs/op

BenchmarkPoolSet-10                     459525678                2.343 ns/op           0 B/op          0 allocs/op
BenchmarkSyncPoolSet-10                 36195546                31.87 ns/op           53 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                19910444                60.20 ns/op           49 B/op          1 allocs/op
BenchmarkSyncPoolAnySet-10              37061773                32.09 ns/op           51 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  44008232                27.56 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              132446389                8.936 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10             68831840                17.20 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           45324494                26.52 ns/op           24 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.320s
```



# []byte-10

```
package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 10

type Model []byte

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		100,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	}}
	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

```



```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     72209403                16.02 ns/op           16 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 19514997                59.89 ns/op           40 B/op          2 allocs/op
BenchmarkPoolsZeroGet-10                14939905                81.41 ns/op           69 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              20102354                61.71 ns/op           40 B/op          2 allocs/op

BenchmarkPoolSet-10                     524742518                2.392 ns/op           0 B/op          0 allocs/op
BenchmarkSyncPoolSet-10                 36350096                32.01 ns/op           52 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                19957645                60.70 ns/op           51 B/op          1 allocs/op
BenchmarkSyncPoolAnySet-10              36793496                31.78 ns/op           50 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  43356190                27.23 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              141722839                8.459 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10             69734329                17.36 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           44950764                27.18 ns/op           24 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.429s
```



# []byte-30

```
package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 30

type Model []byte

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		100,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	}}
	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

```



```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     63825824                18.87 ns/op           32 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 19097864                63.34 ns/op           56 B/op          2 allocs/op
BenchmarkPoolsZeroGet-10                13771779                85.68 ns/op           75 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              18747790                64.01 ns/op           56 B/op          2 allocs/op

BenchmarkPoolSet-10                     491328310                2.307 ns/op           0 B/op          0 allocs/op
BenchmarkSyncPoolSet-10                 37005722                31.32 ns/op           53 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                19661788                60.53 ns/op           50 B/op          1 allocs/op
BenchmarkSyncPoolAnySet-10              36933588                32.39 ns/op           51 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  42862244                27.56 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              131738504                8.456 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10             68141973                17.25 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           45427515                27.02 ns/op           24 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.284s
```



# []byte-100

```
package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 100

type Model []byte

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		100,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return make(Model, 0, Size)
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	}}
	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return make(Model, 0, Size)
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
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return make(Model, 0, Size)
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

```



```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     47410918                24.19 ns/op          112 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 17007213                69.08 ns/op          136 B/op          2 allocs/op
BenchmarkPoolsZeroGet-10                12689382                94.21 ns/op          162 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              17078760                69.72 ns/op          136 B/op          2 allocs/op

BenchmarkPoolSet-10                     530508962                2.278 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroSet-10                19863081                59.11 ns/op           54 B/op          1 allocs/op
BenchmarkSyncPoolSet-10                 37128951                32.47 ns/op           50 B/op          1 allocs/op
BenchmarkSyncPoolAnySet-10              36999969                30.97 ns/op           54 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  44112104                27.17 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              142212831                8.458 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10             69485826                17.34 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           42334699                26.30 ns/op           24 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.201s
```



# struct

```
package pkg

import (
	"github.com/colega/zeropool"
	"sync"
	"testing"
)

var Size = 100

type Model struct {
	Name    string
	Age     string
	Sex     int
	Address map[string]string
}

func BenchmarkPoolGet(b *testing.B) {
	var p = NewPool[Model](10,
		100,
		func(size int) Model {
			return Model{}
		})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return Model{}
		})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolGetSet(b *testing.B) {
	var p = NewPool[Model](10,
		10,
		func(size int) Model {
			return Model{}
		})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)

	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return Model{}
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
		return Model{}
	}}
	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolGetSet(b *testing.B) {
	var p = &sync.Pool{New: func() any {
		return Model{}
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
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkPoolsZeroSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkPoolsZeroGetSet(b *testing.B) {
	var p = zeropool.New[Model](func() Model {
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

func BenchmarkSyncPoolAnyGet(b *testing.B) {
	var p = New[Model](func() Model {
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}
}

func BenchmarkSyncPoolAnySet(b *testing.B) {
	var p = New[Model](func() Model {
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		p.Put(Model{})
	}
}

func BenchmarkSyncPoolAnyGetSet(b *testing.B) {
	var p = New[Model](func() Model {
		return Model{}
	})

	for i := 0; i < b.N; i++ {
		res := p.Get()
		p.Put(res)
	}
}

```



```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     165161866                7.291 ns/op           0 B/op          0 allocs/op
BenchmarkSyncPoolGet-10                 26614545                45.27 ns/op           48 B/op          1 allocs/op
BenchmarkPoolsZeroGet-10                18936205                60.88 ns/op           73 B/op          1 allocs/op
BenchmarkSyncPoolAnyGet-10              24199503                48.49 ns/op           48 B/op          1 allocs/op

BenchmarkPoolSet-10                     511545746                2.465 ns/op           0 B/op          0 allocs/op
BenchmarkSyncPoolSet-10                 32130700                36.60 ns/op           75 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                20693980                56.30 ns/op           75 B/op          1 allocs/op
BenchmarkSyncPoolAnySet-10              31029834                36.56 ns/op           76 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  37098963                31.77 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              135724682                8.484 ns/op           0 B/op          0 allocs/op
BenchmarkPoolsZeroGetSet-10             54887565                21.87 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           38373853                31.53 ns/op           48 B/op          1 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.981s

```



# *struct

```
$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     49017186                23.82 ns/op           48 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 26738463                43.16 ns/op           48 B/op          1 allocs/op
BenchmarkPoolsZeroGet-10                16133793                76.00 ns/op           82 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              27807976                43.58 ns/op           48 B/op          1 allocs/op

BenchmarkPoolSet-10                     54906189                22.55 ns/op           48 B/op          1 allocs/op
BenchmarkSyncPoolSet-10                 35368321                31.38 ns/op           77 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                15691855                71.98 ns/op           85 B/op          2 allocs/op
BenchmarkSyncPoolAnySet-10              36251221                30.95 ns/op           77 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  44175098                28.82 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              28541131                43.91 ns/op           48 B/op          1 allocs/op
BenchmarkPoolsZeroGetSet-10             67134586                17.21 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           140934553                8.447 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.022s


$ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/data-any-pool/pkg
BenchmarkPoolGet-10                     49546213                23.92 ns/op           48 B/op          1 allocs/op
BenchmarkSyncPoolGet-10                 28634264                42.25 ns/op           48 B/op          1 allocs/op
BenchmarkPoolsZeroGet-10                15715891                77.34 ns/op           82 B/op          2 allocs/op
BenchmarkSyncPoolAnyGet-10              26093337                42.90 ns/op           48 B/op          1 allocs/op

BenchmarkPoolSet-10                     52100950                22.01 ns/op           48 B/op          1 allocs/op
BenchmarkSyncPoolSet-10                 36083352                30.83 ns/op           76 B/op          1 allocs/op
BenchmarkPoolsZeroSet-10                16310265                70.78 ns/op           84 B/op          2 allocs/op
BenchmarkSyncPoolAnySet-10              36122137                30.95 ns/op           77 B/op          1 allocs/op

BenchmarkPoolGetSet-10                  44636697                27.17 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolGetSet-10              28525357                42.07 ns/op           48 B/op          1 allocs/op
BenchmarkPoolsZeroGetSet-10             70853936                16.97 ns/op            0 B/op          0 allocs/op
BenchmarkSyncPoolAnyGetSet-10           142322916                8.439 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/zqlpaopao/tool/data-any-pool/pkg     16.860s

```

