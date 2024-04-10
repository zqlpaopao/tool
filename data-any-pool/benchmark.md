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

