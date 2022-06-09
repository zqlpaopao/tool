package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"net"
	"runtime"

	"net/http"
	"time"
)

//https://beego.vip/docs/module/httplib.md
func main() {
	var tp http.RoundTripper = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		Dial:                   nil,
		DialTLSContext:         nil,
		DialTLS:                nil,
		TLSClientConfig:        nil,
		TLSHandshakeTimeout:    0,
		DisableCompression:     false,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,//可选）指定要发送到的标头 连接请求期间的代理。 要动态设置标头，请参阅GetProxyConnectHeader。
		GetProxyConnectHeader:  nil,//（可选）指定要返回的func 在向的连接请求期间发送到代理URL的标头 ip：端口目标。 如果返回错误，则传输的往返将失败 那个错误。它可以返回（nil，nil）以不添加头。 如果GetProxyConnectHeader为非nil，则ProxyConnectHeader为 忽略。
		MaxResponseHeaderBytes: 0,//指定数量限制服务器响应中允许使用响应字节 标题。 零表示使用默认限制。4 << 10
		WriteBufferSize:        0,//指定所用写缓冲区的大小 向运输部门写信时。 如果为零，则使用默认值（当前为4KB）。4 << 10
		ReadBufferSize:         0,//ReadBufferSize指定所用读取缓冲区的大小 从传输读取时。 如果为零，则使用默认值（当前为4KB）。
		ForceAttemptHTTP2:      false,//ForceAttemptHTTP2控制当 提供了Dial、DialTLS或DialContext func或TLSClientConfig。 默认情况下，使用任何这些字段都会保守地禁用HTTP/2。 使用自定义拨号程序或TLS配置并仍尝试HTTP/2 升级，将其设置为true。
		ExpectContinueTimeout: 1 * time.Second,//如果非零，则指定完全恢复后等待服务器第一个响应标头的时间 如果请求具有 “Expect:100 continue”标题。零表示没有超时，并且 使正文立即发送，无需 正在等待服务器批准。 此时间不包括发送请求标头的时间。
		DisableKeepAlives: true,
		ResponseHeaderTimeout: 10*time.Second,//如果非零，则指定完全关闭后等待服务器响应标头的时间 编写请求（包括其正文，如果有）。这时间不包括读取响应正文的时间。

		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,

		MaxIdleConns:          50,
		MaxIdleConnsPerHost :50,//如果非零）控制最大空闲 （保持活动）每个主机要保持的连接。如果为零， 使用DefaultMaxIdleConnsPerHost=2。
		MaxConnsPerHost: 10,//MaxConnsPerHost可以选择限制 每个主机的连接数，包括拨号中的连接数， 活动和空闲状态。违反限制时，拨号将被阻止。 零意味着没有限制。
		IdleConnTimeout:       90 * time.Second,////是空闲的最长时间 保持活动状态）连接在关闭前将保持空闲状态 它本身  零意味着没有限制

		//		MaxIdleConns:          100,
		//		IdleConnTimeout:       90 * time.Second,
		//		TLSHandshakeTimeout:   10 * time.Second,
		//		ExpectContinueTimeout: 1 * time.Second,
		//		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	req := httplib.Post("http://beego.vip/")
	req.SetTransport(tp)


	for i:= 0; i< 3;i++{
		r := 	req.GetRequest()
		fmt.Printf("%#v\n",r)
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
	req.Bytes()
}

