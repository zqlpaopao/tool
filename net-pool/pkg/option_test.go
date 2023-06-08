package pkg

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func BenchmarkInit(b *testing.B) {
	//go server()
}

func BenchmarkNewPoolWithConfig(b *testing.B) {
	p := NewPoolWithConfig[net.Conn](&Config[net.Conn]{
		InitialCap: 3,
		MaxCap:     5,
		MaxIdle:    1,
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
	})
	if p.Error() != nil {
		fmt.Println("err=", p.err)
	}
	p.InitConn()

	for i := 0; i < 1000; i++ {
		//从连接池中取得一个连接
		v, err := p.Get()
		if err != nil {
			fmt.Println("p.Get", err)
		}
		_, errs := v.Conn.Write([]byte("hello world"))
		if errs != nil {
			fmt.Println("v.Conn.Write", errs)
		}
		//将连接放回连接池中
		p.Put(v)

		//查看当前连接中的数量
		//current := p.Len()
		//fmt.Println("len=", current)
	}

}

func BenchmarkDefault(b *testing.B) {
	for i := 0; i < 900; i++ {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println(err)
		}
		_, errs := c.Write([]byte("hello world"))
		if errs != nil {
			fmt.Println("v.Conn.Write", errs)
		}
	}
}

const addr string = "127.0.0.1:8080"

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		_, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		//fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		//go handleRequest(conn)
		//fmt.Println("server----for")
	}
}

func handleRequest(conn net.Conn) {
	for {
		var b = make([]byte, 100)
		_, err := conn.Read(b)
		if err != nil {
			fmt.Println("handleRequest", conn.RemoteAddr(), err)
		}
		fmt.Println("handleRequest", string(b), conn.RemoteAddr())
	}

}
