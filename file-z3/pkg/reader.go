package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Resp struct {
	StartTime time.Time
	EndTime   time.Time
}

type ReaderMan struct {
	fileSize int64
	readSize int
	opt      *Option
	errs     []error
	fd       *os.File
	res      *Resp
	File
}

func NewReaderMan(f ...OptionFunc) *ReaderMan {
	return &ReaderMan{
		opt:  NewOption(f...),
		File: NewFileOp(),
		res:  &Resp{},
	}
}

// ParamsCheck params check
func (r *ReaderMan) ParamsCheck() {
	if r.opt.filePath == "" {
		r.errs = append(r.errs, FileEmtErr)
		return
	}
	if r.opt.tidyData == nil {
		r.errs = append(r.errs, TidyDataEmtErr)
		return
	}
}

// Init init params
func (r *ReaderMan) Init() {
	if len(r.errs) > 0 {
		return
	}
	var err error
	if r.fd, err = os.Open(r.opt.filePath); nil != err {
		r.errs = append(r.errs, err)
		return
	}

	var fileState os.FileInfo
	if fileState, err = r.fd.Stat(); err != nil {
		r.errs = append(r.errs, err)
		return
	}

	r.res.StartTime, r.fileSize = time.Now(), fileState.Size()
}

func (r *ReaderMan) Doing() {
	var (
		wg = &sync.WaitGroup{}
		w1 = &sync.WaitGroup{}
	)
	wg.Add(r.opt.workerNum)
	for i := 0; i < r.opt.workerNum; i++ {
		go r.StartWorker(wg)
	}
	w1.Add(1)
	go r.Read(w1)
	w1.Wait()
	close(r.opt.dataChan)
	wg.Wait()
	r.End()
}

// Read start read file
func (r *ReaderMan) Read(wg *sync.WaitGroup) {
	if len(r.errs) > 0 {
		return
	}
	var rd *bufio.Reader

	rd = bufio.NewReader(r.fd)

	for {
		buf := r.opt.GetByteSliceBuf()
		n, err := rd.Read(buf)
		r.readSize += n
		buf = buf[:n]

		if n == 0 {
			if err != nil && err != io.EOF {
				fmt.Println(err)
				r.errs = append(r.errs, err)
				continue
			}
			if err == io.EOF {
				break
			}
			return
		}

		endData, err := rd.ReadBytes(r.opt.end)
		if err != io.EOF {
			buf = append(buf, endData...)
		}

		if r.opt.checkData != nil && r.opt.checkData(buf) {
			r.opt.dataChan <- buf
			r.opt.PutByteSliceBuf(buf)
			continue
		}
		r.opt.dataChan <- buf
		r.opt.PutByteSliceBuf(buf)

	}
	wg.Done()
}

// StartWorker start customer data
func (r *ReaderMan) StartWorker(wg *sync.WaitGroup) {
	defer r.opt.panicSave()
	for {
		select {
		case v, ok := <-r.opt.dataChan:
			if !ok {
				goto END
			}
			r.opt.tidyData(v)
		}
	}
END:
	wg.Done()
}

func (r *ReaderMan) End() {
	r.res.EndTime = time.Now()
}

// Error errors
func (r *ReaderMan) Error() []error {
	return r.errs
}

func (r *ReaderMan) Code() int {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderMan) GetResp() *Resp {
	return r.res
}
