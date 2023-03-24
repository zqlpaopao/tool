package pkg

import "os"

type FileOp struct {
}

func NewFileOp() *FileOp {
	return &FileOp{}
}

func (f *FileOp) Open(s string) (*os.File, error) {
	return os.Open(s)
}

func (f *FileOp) Close() error {
	//TODO implement me
	panic("implement me")
}
