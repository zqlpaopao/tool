package main

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/zqlpaopao/tool/http_page_request/pkg"
	"os"
	"sync"
)

type Datas struct {
	lock *sync.Mutex
	data map[string]string
}

func main() {
	var data = &Datas{
		lock: &sync.Mutex{},
		data: make(map[string]string, 10000),
	}
	handler := pkg.NewPageHandlerWithOptions[Req, *Datas](
		Req{
			Page: Page{
				Start: 0,
				Limit: 1,
			},
			Con: Condition{},
		},
		resFn,
		totalFn,
		data,
		pkg.WithUrl("xxxxxxx.com"),
	)
	err := handler.DO()
	fmt.Println(err)
	fmt.Println(data.data)

}

// 获取共有多少数据
func totalFn(url string, param Req) int {
	var res *Resp
	g, err := grequests.Post(url, &grequests.RequestOptions{
		JSON: param,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = g.JSON(&res); nil != err {
		fmt.Println(err)
		os.Exit(2)
	}
	return res.Data.Count
}

// 每个数据处理
func resFn(url string, p Req, page, limit int, resd *Datas) *pkg.ErrInfo[Req] {

	p.Page.Start, p.Page.Limit = page, limit
	fmt.Println(p)

	var res *Resp
	g, err := grequests.Post(url, &grequests.RequestOptions{
		JSON: p,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	if err = g.JSON(&res); nil != err {
		fmt.Println(err)
		os.Exit(3)
	}

	for _, v := range res.Data.Info {
		resd.lock.Lock()
		resd.data[v.Code] = v.Code
		resd.lock.Unlock()

	}

	return nil
}

type Req struct {
	Con  Condition `json:"condition"`
	Page Page      `json:"page"`
}

type Page struct {
	Start int `json:"start"`
	Limit int `json:"limit"`
}

type Condition struct {
}

type Resp struct {
	ErrorMsg  string `json:"error_msg"`
	Data      Data   `json:"data"`
	ErrorCode int    `json:"error_code"`
	Result    bool   `json:"result"`
}

type Data struct {
	Info  []Info `json:"info"`
	Count int    `json:"count"`
}

type Info struct {
	Code string `json:"code"`
}
