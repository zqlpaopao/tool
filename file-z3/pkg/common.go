package pkg

import "errors"

const (
	ByteDataSize  = 0x51199
	CacheByteChSi = 0x199999
	DataChanSize  = 0x999399
	WorkerBUm     = 0x100
	End           = 'n'
)

var (
	FileEmtErr     = errors.New("file is empty")
	TidyDataEmtErr = errors.New("tidy data function is empty")
)
