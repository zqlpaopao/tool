package pkg

import "errors"

const (
	ByteDataSize  = 0x51199
	CacheByteChSi = 0x199999
	Customer      = 100
	ReaderSize    = 0x49999
	DataChSize    = 50000
	WorkerBUm     = 0x9
	ReadWorkerNum = 100
	End           = '\n'
)

var (
	FileEmtErr     = errors.New("file is empty")
	TidyDataEmtErr = errors.New("tidy data function is empty")
)
