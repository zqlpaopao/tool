package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_ping_agent/pkg"
	"github.com/zqlpaopao/tool/string-byte/src"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	go Loop()
	http.ListenAndServe("0.0.0.0:6061", nil)

}

func Loop() {
	var pool = ping.NewPoolData(10)
	ping.InitUUIDPool(60, 15)
	p := ping.NewPoolWithOptions(
		//ping.WithProtocol("icmp"),
		ping.WithProtocol("icmp"),
		ping.WithOnRevFunc(func(pings *ping.Ping, packet *ping.Packet) {
			//fmt.Println()
			//fmt.Printf("%#v\n", packet)
			st := ping.StatisticsLog(pings, packet)
			ping.Debug(st)
			fmt.Println("uuid1--pings", pings.Uuid().String())
			fmt.Println("uuid1--packs", packet.Uuid())

			//回收资源
			ping.SetUUID(*pings.Uuid())
			var array [36]byte
			copy(array[:], src.String2Bytes(packet.Uuid())[0:36])
			ping.SetByte(array)

			//最后，不然有数据race
			pool.SetPing(pings)

		}),
		ping.WithErrCallBack(func(ping *ping.Ping, err error) {
			if err == nil {
				return
			}

			if ping == nil {
				return
			}
			fmt.Printf("%#v", ping)
			fmt.Println(ping.Addr(), err)
			os.Exit(1)
		}),
	).Run()
	err := p.Error()
	if err != nil {
		panic(err)
	}
	for {

		IPV4Addr(pool, p, "110.242.68.66", 110, 242, 68, 66)
		IPV4Addr(pool, p, "123.151.137.18", 123, 151, 137, 18)
		IPV4Addr(pool, p, "203.205.254.157", 203, 205, 254, 157)


		time.Sleep(100 * time.Millisecond)
	}
}

func IPV4Addr(pool *ping.PoolData, p *ping.Pool, addr string, a, b, c, d byte) {
	var err error
	bd := pool.GetPing()
	u := ping.NewUUid()

	bd.SetDstAddr(addr).SetIpV4().SetResolveIpAddr(&net.IPAddr{IP: net.IPv4(a, b, c, d), Zone: ""}).SetSize(60).SetTtl(60).SetUUid(&u)

	if err = p.Submit(bd); nil != err {
		panic(err)
	}

}

func hostIpV4(pool *ping.PoolData, p *ping.Pool) {
	//var err error
	//qq := pool.GetPing()
	//if err = qq.SetAddr("qq.com"); nil != err {
	//	panic(err)
	//}
	//qq.SetSize(40).SetTtl(60).SetUUid(uuid.New())
	//
	//if err = p.Submit(qq); nil != err {
	//	panic(err)
	//}

}

//
//use std::{net::{SocketAddr}, sync::Arc, time::Duration};
//use bytes::{Bytes, BytesMut};
//use anyhow::{Result};
//use async_channel::{bounded, Receiver, Sender};
//use tokio::{signal, net::UdpSocket, time::sleep};
//
//
///// worker 线程
//async fn worker(worker_name: String, req_receiver: Arc<Receiver<(SocketAddr, Bytes)>>, res_sender: Arc<Sender<(SocketAddr, Bytes)>>) -> Result<()>{
//println!("{} running...", worker_name);
//loop {
//tokio::select! {
//res = req_receiver.recv() => {
//let (src, msg) = match res {
//Ok(r) => r,
//Err(e) => {
//println!("{}, Unable to read data, error:{}", worker_name, e);
//return Ok(());
//}
//};
//
//// 处理逻辑, 等待1s...
//sleep(Duration::from_secs(1)).await;
//println!("{}, process the request, src:{:?}", worker_name, src);
//
//// 结果发送到回复队列
//match res_sender.send((src, msg)).await {
//Ok(_) => {},
//Err(e) => println!("{}, Failed to send message to reply queue, error:{}", worker_name, e),
//}
//}
//_ = signal::ctrl_c() => {
//println!("{}, exit...", worker_name);
//req_receiver.close();
//return Ok(());
//}
//}
//}
//
//}
//
//// server线程
//async fn server(socket: UdpSocket, req_sender: Arc<Sender<(SocketAddr, Bytes)>>, res_receiver: Arc<Receiver<(SocketAddr, Bytes)>>) -> Result<()> {
//println!("server running...");
//loop {
//let mut buf = BytesMut::with_capacity(1024);
//buf.resize(1024, 0);
//
//tokio::select! {
//
//// socket中读取到数据, 发送到请求队列.
//res = socket.recv_from(&mut buf) => {
//let (len, src) = match res {
//Ok(r) => r,
//Err(e) => {
//println!("Unable to read data, error:{}", e);
//continue;
//}
//};
//println!("receive socket data, send to request queue. src:{:?}...", src);
//buf.resize(len, 0);
//req_sender.send((src, buf.freeze())).await?;
//}
//
//// 接收回复队列的数据发送到socket客户端
//Ok((src, msg)) = res_receiver.recv() => {
//println!("receive response queue data, send to socket client:{:?}...", src);
//match socket.send_to(&msg, src).await {
//Ok(_) => continue,
//Err(e) => println!("Failed to send data to client({:?}), error:{:?}", src, e),
//}
//}
//
//_ = signal::ctrl_c() => {
//println!("exit...");
//res_receiver.close();
//return Ok(());
//}
//
//}
//}
//}
//
//
//#[tokio::main]
//async fn main() -> Result<()> {
//
//let addr = "0.0.0.0:3053".parse::<SocketAddr>()?;
//let socket = UdpSocket::bind(addr).await?;
//
//// 请求队列
//let (req_sender, req_receiver) = bounded::<(SocketAddr, Bytes)>(1024);
//// 回复队列
//let (res_sender, res_receiver) = bounded::<(SocketAddr, Bytes)>(1024);
//
//let req_sender = Arc::new(req_sender);
//let req_receiver = Arc::new(req_receiver);
//
//let res_sender = Arc::new(res_sender);
//let res_receiver = Arc::new(res_receiver);
//
//for i in 0..4 {
//let req_receiver = req_receiver.clone();
//let res_sender = res_sender.clone();
//// worker 线程
//tokio::spawn(async move {
//let worker_name = format!("worker:<{:?}>", i);
//let _ = worker(worker_name, req_receiver, res_sender).await;
//});
//}
//
//// server 线程
//server(socket, req_sender, res_receiver).await
//}
//
//
