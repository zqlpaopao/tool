package src

const (
	syncGoNum       = 10
	syncBuffSize    = 1024
	maxSyncBuffSize = 67021478
	bufferMax       = 500
)

type CallerInfo struct {
	FileLine int
	FuncName string
	FilePath string
}
