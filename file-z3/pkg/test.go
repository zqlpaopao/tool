package pkg

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const filePath = "xxx/log.log"

func Test() {
	fileName := "./test.log"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Could not able to get the file stat")
		return
	}
	loop := fileStat.Size() / 30

	for i := 0; i <= int(loop); i++ {
		bytes := make([]byte, 30)
		if _, err := file.ReadAt(bytes, 30*int64(i)); nil != err {
			fmt.Println(err, "-----", i)
		}
		fmt.Println("i-----", i)
		fmt.Println(string(bytes))
	}

	fmt.Println(fileStat.Size())
}

func Old() {
	s := time.Now()
	fileName := filePath
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return
	}

	defer file.Close() //close after checking err

	//queryStartTime, err := time.Parse("2006-01-02T15:04:05.0000Z", startTimeArg)
	//if err != nil {
	//	fmt.Println("Could not able to parse the start time", startTimeArg)
	//	return
	//}
	//
	//queryFinishTime, err := time.Parse("2006-01-02T15:04:05.0000Z", finishTimeArg)
	//if err != nil {
	//	fmt.Println("Could not able to parse the finish time", finishTimeArg)
	//	return
	//}

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Could not able to get the file stat")
		return
	}

	fileSize := fileStat.Size()
	offset := fileSize - 1
	lastLineSize := 0

	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "n" {
			break
		}
		offset--
		lastLineSize += n
	}

	lastLine := make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset+1)

	if err != nil {
		fmt.Println("Could not able to read last line with offset", offset, "and lastline size", lastLineSize)
		return
	}
	fmt.Println("fileSize := fileStat.Size()", fileSize)

	//logSlice := strings.SplitN(string(lastLine), ",", 2)
	//logCreationTimeString := logSlice[0]

	//lastLogCreationTime, err := time.Parse("2006-01-02T15:04:05.0000Z", logCreationTimeString)
	//if err != nil {
	//	fmt.Println("can not able to parse time : ", err)
	//}

	//if lastLogCreationTime.After(queryStartTime) && lastLogCreationTime.Before(queryFinishTime) {
	Process(file)
	//}

	fmt.Println("nTime taken - ", time.Since(s))
}

func Process(f *os.File) error {

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}
	//250000

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)
	var in = 0

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]
		in += n
		fmt.Println("len(buf),n", len(buf), n)
		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		wg.Add(1)
		go func() {
			len := len(buf)
			len = len
			stringPool = stringPool
			//ProcessChunk(buf, &linesPool, &stringPool)
			wg.Done()
		}()

	}

	wg.Wait()
	fmt.Println(count)
	fmt.Println("count", in)
	return nil
}

var count int64

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {

	var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}
	atomic.AddInt64(&count, int64(noOfThread))
	//fmt.Println("noOfThread", noOfThread)
	for i := 0; i < noOfThread; i++ {

		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done() //to avaoid deadlocks
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				//logSlice := strings.SplitN(text, ",", 2)
				//logCreationTimeString := logSlice[0]

				//logCreationTime, err := time.Parse("2006-01-02T15:04:05.0000Z", logCreationTimeString)
				//if err != nil {
				//	fmt.Printf("n Could not able to parse the time :%s for log : %v", logCreationTimeString, text)
				//	return
				//}

				//if logCreationTime.After(start) && logCreationTime.Before(end) {
				//fmt.Println(text)
				//}
			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	wg2.Wait()
	logsSlice = nil
}

func MuiltGo() {

	t := time.Now()
	fileName := filePath
	desfileName := filePath
	sfile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(nil)
	}

	info, _ := os.Stat(fileName)
	size := info.Size()

	loop := size / 300000
	//fmt.Println(loop)
	//os.Exit(1)
	//var scount int64 = 1
	//if size%2 == 0 {
	//	scount *= 2
	//} else if size%3 == 0 {
	//	scount *= 3
	//} else {
	//	scount *= 1
	//}

	desF, err := os.Create(desfileName)
	if err != nil {
		fmt.Println(err)
	}

	var num chan int = make(chan int, 1000)

	wg1 := sync.WaitGroup{}
	wg1.Add(100)
	for i := 0; i < 100; i++ {
		go cum(sfile, num, &wg1, desF)
	}

	for i := 0; i <= int(loop); i++ {
		num <- i
	}
	close(num)
	wg1.Wait()

	//
	//wg := sync.WaitGroup{}
	//wg.Add(int(si))
	//for i := 0; i < int(si); i++ {
	//	go func(vs int) {
	//		//申明一个byte
	//
	//		//从指定位置开始写
	//		desF.WriteAt(b, int64(vs)*si)
	//		//从指定位置开始写
	//		wg.Done()
	//
	//	}(i)
	//}
	//wg.Wait()
	fmt.Println(time.Now().Sub(t))
	defer sfile.Close()
	defer desF.Close()
}

func cum(file *os.File, num chan int, wg *sync.WaitGroup, desfileName *os.File) {
	for {
		select {
		case v, ok := <-num:
			if !ok {
				goto END
			}
			b := make([]byte, 300000)
			//从指定位置开始读
			if _, err := file.ReadAt(b, int64(v)*300000); nil != err {
				fmt.Println(err)
			}
			//fmt.Println(string(b))
			////从指定位置开始写
			//if _, err := desfileName.WriteAt(b, int64(v)*300000); nil != err {
			//	fmt.Println(err)
			//	os.Exit(1)
			//}
		}
	}
END:
	wg.Done()
}
