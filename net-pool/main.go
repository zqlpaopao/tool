package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/net-pool/pkg"
	"net"
	"os"
	"time"
)

// https://github.com/silenceper/pool/tree/master
// https://github.com/fatih/pool
/*
	//支持grpc的连接池和http的连接池
*/

const addr string = "127.0.0.1:8080"

func main() {

	go server()
	//等待tcp server启动
	time.Sleep(2 * time.Second)
	client()
	select {}
}

func client() {

	p := pkg.NewPoolWithConfig[net.Conn](&pkg.Config[net.Conn]{
		InitialCap:    3,
		MaxCap:        5,
		MaxIdle:       1,
		IsCheck:       true,
		CheckInterval: 2 * time.Second,
		Factory: func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		},
		Close: func(conn net.Conn) {
			_ = conn.Close()
		},
		Ping: func(conn net.Conn) error {
			return nil
		},
		IdleTimeout: 15 * time.Second,
	}).InitConn()

	err := p.Error()
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 10; i++ {
		go func() {
			for {
				//从连接池中取得一个连接
				v, err := p.Get()
				if err != nil {
					fmt.Println("p.Get", err)
				}
				//fmt.Println("get-------------------")
				_, errs := v.Conn.Write([]byte("hello world"))
				if errs != nil {
					//fmt.Println("v.Conn.Write", errs)
				}
				//fmt.Println("write----n", n)
				//将连接放回连接池中
				p.Put(v)
				//fmt.Println("put-------------------")

				//查看当前连接中的数量
				//current := p.Len()
				//fmt.Println("len=", current)
				//time.Sleep(1 * time.Second)
			}
		}()

	}

	time.Sleep(10 * time.Second)
	p.Release()

}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
		fmt.Println("server----for")
	}
}

func handleRequest(conn net.Conn) {
	for {
		var b = make([]byte, 100)
		_, err := conn.Read(b)
		if err != nil {
			//fmt.Println("handleRequest", conn.RemoteAddr(), err)
		}
		//fmt.Println("handleRequest", string(b), conn.RemoteAddr())
	}

}
