package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/zqlpaopao/tool/http/src"
)

//https://beego.vip/docs/module/httplib.md
func main() {
	libHttp()
}

func libHttp() {
	c, err := httplib.Get("https://baidu.com").Debug(true).Response()
	fmt.Println(c)
	fmt.Println(err)
}

func YuanSheng() {
	c, err := src.CurlCtx.NewRequest(src.RequestArgs{
		RequestUrl: "https://baidu.com",
		RequestParams: map[string]string{
			"a": "b",
		},
		RequestType:   src.RequestPost,
		RequestHeader: src.RequestForm.String(),
		SetTimeOut:    10,
	})

	fmt.Println(err)

	r, err := c.Send()
	fmt.Println(r)
	fmt.Println(err)
}
