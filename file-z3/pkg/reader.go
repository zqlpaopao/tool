package pkg

import (
	"bufio"
	"errors"
	"github.com/sourcegraph/conc/stream"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Resp struct {
	StartTime time.Time
	EndTime   time.Time
	ReadTime  time.Duration
	FileSize  int64
	ReadSize  int64
}

type ReaderMan struct {
	Opt  *Option
	fd   *os.File
	Res  *Resp
	lock *sync.Mutex
	errs []error
}

func NewReaderMan() *ReaderMan {
	return &ReaderMan{
		Opt: &Option{
			IsNewLine:     false,
			ByteDataSize:  ByteDataSize,
			Customer:      Customer,
			CacheByteChSi: CacheByteChSi,
			WorkerNum:     WorkerBUm,
			ReadWorkerNum: ReadWorkerNum,
			ReaderSize:    ReaderSize,
			DataChSize:    DataChSize,
			FilePath:      "",
			byteChan:      nil,
			dataChan:      nil,
			End:           End,
			TidyData:      nil,
			CheckData:     nil,
			PanicSave:     DefaultSavePanic,
		},
		Res: &Resp{},
	}
}

func (r *ReaderMan) Do() *ReaderMan {
	defer r.Opt.PanicSave()
	r.ParamsCheck()
	r.Init()
	r.Doing()
	return r
}

// ParamsCheck params check
func (r *ReaderMan) ParamsCheck() {
	if r.Opt.FilePath == "" {
		r.errs = append(r.errs, FileEmtErr)
		return
	}
	if r.Opt.TidyData == nil {
		r.errs = append(r.errs, TidyDataEmtErr)
		return
	}
}

// Init init params
func (r *ReaderMan) Init() {
	if len(r.errs) > 0 {
		return
	}
	var (
		err       error
		fileState os.FileInfo
	)
	if r.fd, err = os.Open(r.Opt.FilePath); nil != err {
		r.errs = append(r.errs, err)
		return
	}

	if fileState, err = r.fd.Stat(); err != nil {
		r.errs = append(r.errs, err)
		return
	}

	r.Res.StartTime,
		r.Res.FileSize,
		r.Opt.byteChan,
		r.lock,
		r.Opt.curCustomer =
		time.Now(),
		fileState.Size(),
		make(chan *[]byte, r.Opt.CacheByteChSi),
		&sync.Mutex{},
		0

	if !r.Opt.IsNewLine {
		r.Opt.dataChan = make(chan *[]byte, r.Opt.DataChSize)
	}
	//go func() {
	//	for i := 0; i < r.Opt.CacheByteChSi/2; i++ {
	//		r.Opt.PutByteSliceBuf(&[]byte{})
	//	}
	//}()
}

func (r *ReaderMan) Doing() {
	//go func() {
	//	for {
	//		fmt.Println(runtime.NumGoroutine())
	//		time.Sleep(100 * time.Millisecond)
	//	}
	//
	//}()
	r.Read()
	r.End()
}

// Read start read file
func (r *ReaderMan) Read() {
	if r.Opt.IsNewLine {
		r.ReadNewLine()
		return
	}
	r.ReadStream()
}

func (r *ReaderMan) ReadNewLine() {
	var (
		rd *bufio.Reader
		t  = time.Now()
	)
	rd = bufio.NewReaderSize(r.fd, r.Opt.ReaderSize)

	for {
		//buf := r.Opt.GetByteSliceBuf()
		buf := make([]byte, r.Opt.ByteDataSize)
		n, err := rd.Read(buf)

		r.Res.ReadSize += int64(n)

		if cap(buf) != n {
			buf = buf[:n]
		}

		if n == 0 {
			if err != nil && err != io.EOF {
				r.errs = append(r.errs, errors.New("io.EOF"+err.Error()))
				continue
			}
			if err == io.EOF {
				break
			}
			return
		}

		endData, err := rd.ReadBytes(r.Opt.End)
		if err != io.EOF {
			buf = append(buf, endData...)
		}
		if r.Opt.CheckData != nil && r.Opt.CheckData(&buf) {
			r.Callback(&buf)
			continue
		}
		r.Callback(&buf)

	}
	r.Res.ReadTime = time.Now().Sub(t)
}

func (r *ReaderMan) Callback(b *[]byte) {
	if b == nil {
		return
	}
	for {
		if atomic.LoadInt64(&r.Opt.curCustomer) < 100 {
			go r.runningGo(b)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (r *ReaderMan) runningGo(b *[]byte) {
	defer r.Opt.PanicSave()
	r.Opt.TidyData(b)
	//r.Opt.PutByteSliceBuf(b)
	atomic.AddInt64(&r.Opt.curCustomer, -1)
}

func (r *ReaderMan) ReadStream() {
	var t = time.Now()
	streams := stream.New()
	r.Customer()

	r.StreamProducer(streams)
	streams.Wait()

	r.Res.ReadTime = time.Now().Sub(t)
	close(r.Opt.dataChan)

}

func (r *ReaderMan) StreamProducer(streams *stream.Stream) {
	var loop int

	loop, r.Opt.dataChan = r.makeLoop(), make(chan *[]byte, 3+1)

	for i := 0; i <= loop; i++ {
		num := i
		streams.Go(func() stream.Callback {
			buf := r.ReadStreamByte(num)
			return func() { r.pushData(buf) }
		})
	}
}

func (r *ReaderMan) Customer() {
	if len(r.errs) > 0 {
		return
	}
	if !r.Opt.IsNewLine {
		r.Opt.WorkerNum = 1
	}
	for i := 0; i < r.Opt.WorkerNum; i++ {
		go r.StartWorker()
	}
}

// StartWorker start customer data
func (r *ReaderMan) StartWorker() {
	defer r.Opt.PanicSave()
	for {
		select {
		case v, ok := <-r.Opt.dataChan:
			if !ok {
				goto END
			}
			r.Opt.TidyData(v)
			//r.Opt.PutByteSliceBuf(v)
		default:

		}
	}
END:
}

func (r *ReaderMan) makeLoop() int {
	return int(r.Res.FileSize) / r.Opt.ByteDataSize
}

func (r *ReaderMan) ReadStreamByte(num int) *[]byte {
	defer r.Opt.PanicSave()

	if len(r.errs) > 0 {
		return &[]byte{}
	}

	//buf := r.Opt.GetByteSliceBuf()
	buf := make([]byte, r.Opt.ByteDataSize)

	n, err := r.fd.ReadAt(buf, int64(num*(r.Opt.ByteDataSize)))
	if err != nil {
		r.lock.Lock()
		r.errs = append(r.errs, errors.New("ReadAt"+err.Error()))
		r.lock.Unlock()
	}
	atomic.AddInt64(&r.Res.ReadSize, int64(n))
	buf = buf[:n]

	if n == 0 {
		if err != nil && err != io.EOF {
			r.lock.Lock()
			r.errs = append(r.errs, errors.New("stream io.EOF"+err.Error()))
			r.lock.Unlock()
			return &buf
		}
		if err == io.EOF {
			return &buf

		}
	}
	if r.Opt.CheckData != nil && r.Opt.CheckData(&buf) {
		return &buf
	}
	return &buf
}

func (r *ReaderMan) pushData(buf *[]byte) {
	r.Opt.dataChan <- buf
}

func (r *ReaderMan) End() {
	r.Res.EndTime = time.Now()
	if r.fd != nil {
		if err := r.fd.Close(); nil != err {
			r.errs = append(r.errs, err)
		}
	}
}

// Error errors
func (r *ReaderMan) Error() []error {
	return r.errs
}

func (r *ReaderMan) GetResp() *Resp {
	return r.Res
}
