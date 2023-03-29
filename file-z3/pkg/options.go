package pkg

import (
	"fmt"
	"os"
	"runtime/debug"
)

type Option struct {
	PanicSave     func()
	CheckData     func(*[]byte) bool
	TidyData      func(*[]byte)
	dataChan      chan *[]byte
	byteChan      chan *[]byte
	FilePath      string
	CacheByteChSi int
	curCustomer   int64
	ReaderSize    int
	ReadWorkerNum int
	WorkerNum     int
	Customer      int64
	DataChSize    int
	ByteDataSize  int
	IsNewLine     bool
	End           byte
}

// type OpFunc func(*Option)
//
//	func NewOption(opt ...OptionFunc) *Option {
//		return clone().WithOptions(opt...)
//	}
//
// // apply assignment function entity
//
//	func (o OpFunc) apply(opt *Option) {
//		o(opt)
//	}
//
// // clone  new object
//
//	func clone() *Option {
//		return &Option{
//			//cacheByteChSi: CacheByteChSi,
//			//byteDataSize:  ByteDataSize,
//			//dataChanSize:  DataChanSize,
//			//end:           End,
//			//workerNum:     WorkerBUm,
//			//readWorkerNum: ReadWorkerBUm,
//			//tidyData:      nil,
//			//checkData:     nil,
//			//panicSave:     defaultSavePanic,
//		}
//	}
//
// // WithOptions Execute assignment function entity
//
//	func (o *Option) WithOptions(opt ...OptionFunc) *Option {
//		for _, v := range opt {
//			v.apply(o)
//		}
//		o.initParams()
//		return o
//	}
//
// // initParams Initialization parameters
//
//	func (o *Option) initParams() {
//		o.byteChan, o.dataChan = make(chan *[]byte, o.cacheByteChSi), make(chan *[]byte, o.dataChanSize)
//		go func() {
//			t := time.Now()
//			for i := 0; i < o.byteDataSize/2; i++ {
//				o.PutByteSliceBuf(&[]byte{})
//			}
//			fmt.Println(time.Now().Sub(t))
//		}()
//	}
//
// // WithIsNewLine is stream reader
//
//	func WithIsNewLine(b bool) OpFunc {
//		return func(o *Option) {
//			o.IsNewLine = b
//		}
//	}
//
// // WithFilePath make byteDataSize
//
//	func WithFilePath(path string) OpFunc {
//		return func(o *Option) {
//			o.filePath = path
//		}
//	}
//
// // WithByteDataSize make byteDataSize
//
//	func WithByteDataSize(num int) OpFunc {
//		return func(o *Option) {
//			o.byteDataSize = num
//		}
//	}
//
// // WithDataChanSize How much to make readyDataSize
//
//	func WithDataChanSize(size int) OpFunc {
//		return func(o *Option) {
//			o.dataChanSize = size
//		}
//	}
//
// // WithHandleGoNum Number of goroutine processed default 256
//
//	func WithHandleGoNum(num int) OpFunc {
//		return func(o *Option) {
//			o.workerNum = num
//		}
//	}
//
// // WithReadWorkerNum Number of goroutine processed default 256
//
//	func WithReadWorkerNum(num int) OpFunc {
//		return func(o *Option) {
//			o.readWorkerNum = num
//		}
//	}
//
// // WithHandleEnd read end byte
//
//	func WithHandleEnd(b byte) OpFunc {
//		return func(o *Option) {
//			o.end = b
//		}
//	}
//
// // WithTidyData will doing table name
//
//	func WithTidyData(f func(*[]byte)) OpFunc {
//		return func(o *Option) {
//			o.tidyData = f
//		}
//	}
//
// // WithCheckData check data function
//
//	func WithCheckData(f func(*[]byte) bool) OpFunc {
//		return func(o *Option) {
//			o.checkData = f
//		}
//	}
//
// // WithPanicSave save panic function
//
//	func WithPanicSave(f func()) OpFunc {
//		return func(o *Option) {
//			o.panicSave = f
//		}
//	}
//

// DefaultSavePanic save panic
func DefaultSavePanic() {
	if err := recover(); nil != err {
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
		os.Exit(3)
	}
}
