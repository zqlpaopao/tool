package main

import (
	"fmt"
	pkg2 "github.com/zqlpaopao/tool/net-pool/pkg"
	"net"
	"time"
)

/*
	//支持grpc的连接池和http的连接池
 */

func main(){
	cp ,err := pkg2.NewItemPool(
		//pkg2.WithAddr("127.0.0.1"),
		//pkg2.WithPort("8080"),
		pkg2.WithMakeConn(func() (pkg2.ConnRes, error) {
			return net.Dial("tcp",":8080")
		}),
		pkg2.WithMaxNum(4),
		pkg2.WithTimeout(10*time.Second),
	)

	fmt.Println(err)

	for {
			c1 := cp.Get()


		c1.GetHttpConn().RemoteAddr()

			fmt.Printf("addr %p\n",c1)
			fmt.Println("NumPooled",cp.NumPooled())
			n1,err := c1.GetHttpConn().Write([]byte("hello word"))
			fmt.Println(n1,"n1")

			fmt.Println("写数据",err)
			buf := make([]byte, 1024)
			fmt.Println(2)
			n, err := c1.GetHttpConn().Read(buf)
			fmt.Println(3)
			fmt.Println("读数据",err)
			fmt.Println("conn1 read : ", string(buf[:n]))
			fmt.Println("end")
			cp.Put(c1)

		fmt.Println("NumPooled",cp.NumPooled())
		time.Sleep(5*time.Second)

	}

}
//https://www.jianshu.com/p/43bb39d1d221
//服务端
//package main
//
//import (
//"fmt"
//"net"
//"time"
//
////"io"
//"log"
//)
//
//func handler(conn net.Conn) {
//	recieveBuffer := []byte("return")
//	for {
//		println("Handling connection! ", conn.RemoteAddr().String(), " connected!")
//
//		// how does it know the end of a message vs the start of a new one?
//		messageSize, err := conn.Write(recieveBuffer)
//		if err != nil {
//			return
//		}
//
//		if messageSize > 0 { // update keep alive since we got a message
//			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
//		}
//	}
//
//	//fmt.Printf("%#v",conn.RemoteAddr())
//	//buf := make([]byte, 1024);
//	//n, _ := conn.(net.Conn).Read(buf);
//	//fmt.Println(buf)
//	////n,err := io.Copy(conn, conn);
//	//fmt.Println(n)
//	////fmt.Println(err)
//}
//
//func main() {
//	lis, err := net.Listen("tcp", ":8080");
//	if err != nil {
//		log.Fatal(err);
//	}
//
//	for {
//		fmt.Println(11)
//		conn, err := lis.Accept();
//		if err != nil {
//			fmt.Println(err)
//			continue;
//		}
//		go handler(conn);
//	}
//}

