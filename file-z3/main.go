package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/file-z3/pkg"
	"sync"
)

const filePath = "mmm/log.log"

//const filePath = "xxxx.log"

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		//r := pkg.NewReaderMan()
		// 读取不写入 chan  3.623350125s  写入chan 28.674761875s
		r := pkg.ReaderMan{
			Opt: &pkg.Option{
				IsNewLine:    true,   //true 是循环读取，false是offset读取
				ByteDataSize: 200000, //读取的 []byte 大小
				Customer:     10,     //暂时无用，offset 最大处理协程1 最大的处理协程数  默认100
				//CacheByteChSi: 300000, //暂时无用 此值别给 走新建会快些 涉及chan 操作会慢
				WorkerNum:     10,       //offset 为保持有序性 只能为1 ，此时默认值无用 默认是9
				ReadWorkerNum: 100,      //默认是100 循环读取 允许开启的消费协程的最大数
				ReaderSize:    320000,   //循环读取的buffer大小
				DataChSize:    100000,   //offset读取的存储数据的大小//为保证有序性 ，只能走chan，chan 影响性能较大
				FilePath:      filePath, //文件路径
				End:           '\n',
				TidyData: func(i *[]byte) {
					//fmt.Println("TidyData", len(*i))
				},
				CheckData: func(i *[]byte) bool {
					//fmt.Println("-----0", string(*i))
					return true
				},
				PanicSave: pkg.DefaultSavePanic,
			},
			Res: &pkg.Resp{},
		}
		r.Do()
		fmt.Println(r.Error())
		fmt.Println(r.GetResp().EndTime.Sub(r.GetResp().StartTime))
		fmt.Printf("%#v", r.GetResp())
		wg.Done()
	}()
	wg.Wait()
}

/*
2020-01-31T20:12:38.1234Z, Some Field, Other Field, And so on, Till new line,...n
*/
