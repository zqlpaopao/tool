package main

import (
	"flag"
	"fmt"
	"github.com/zqlpaopao/tool/go-binlog/pkg"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	runtime.GOMAXPROCS(runtime.NumCPU()/4 + 1)
	flag.Parse()
	cfg := &pkg.Config{
		Host: "127.0.0.1",
		Port: 33306,
		User: "root",
		Pass: "123456",
		ServerId: 2,
		LogFile: "mysql-bin.000001",
		Position: 123,
	}
	srv := pkg.NewServer(cfg,
		pkg.WithOptConTime(10),
		pkg.WithOptPrint(true),
		pkg.WithOptTryDump(true),
		pkg.WithOptTryNum(3))
	go srv.Run()

	//time.Sleep(5*time.Second)
	//fmt.Println(srv.Error())

	select {
	case n := <-sc:
		srv.Quit()
		fmt.Printf("receive signal %v, closing", n)
	}
}
