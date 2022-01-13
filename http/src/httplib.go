package src

import (
	"github.com/beego/beego/v2/client/httplib"
	"sync"
	"time"
)
// https://beego.vip/docs/module/httplib.md
type curl struct {
	sync.Once
	requestTimeout  time.Duration // 请求超时时间
	readDataTimeout time.Duration // 数据读取超时
}

var CurlIns = new(curl)

func init() {
	CurlIns.Once.Do(func() {
		CurlIns.requestTimeout = 3 * time.Second
		CurlIns.readDataTimeout = 2 * time.Second

	})
}

//GetCurlIns  Get instances separately
func (c *curl)GetCurlIns()*curl{
	return &curl{}
}

//SetConnectTimeout Set connect timeout
func (c *curl)SetConnectTimeout(t time.Duration){
	c.requestTimeout = t
}

//SetReadWriteTimeout Set read timeout
func (c *curl)SetReadWriteTimeout(t time.Duration){
	c.readDataTimeout = t
}

func (c *curl) Get(url string, params map[string]string) (rep []byte, err error) {
	ins := httplib.Get(url)
	if len(params) > 0 {
		for k, v := range params {
			ins.Param(k, v)
		}
	}
	if rep, err = ins.SetTimeout(c.requestTimeout, c.readDataTimeout).Bytes(); err != nil {
	}
	return rep, err
}

func (c *curl) PostForm(url string, params map[string]string) (rep []byte, err error) {
	ins := httplib.Post(url)
	if len(params) > 0 {
		for k, v := range params {
			ins.Param(k, v)
		}
	}
	if rep, err = ins.SetTimeout(c.requestTimeout, c.readDataTimeout).Bytes(); err != nil {
	}
	return rep, err
}

func (c *curl) PostJson(url string, params map[string]interface{}) (rep []byte, err error) {
	ins := httplib.Post(url)
	if len(params) > 0 {
		_, err = ins.JSONBody(params)
	}
	if rep, err = ins.SetTimeout(c.requestTimeout, c.readDataTimeout).Bytes(); err != nil {
	}
	return rep, err
}

//PostFile Support file upload
func (c *curl) PostFile(url string, params map[string]string,fileName map[string]string) (rep []byte, err error) {
	ins := httplib.Post(url)
	for k, v := range params {
		ins.Param(k, v)
	}

	for k ,v := range fileName{
		ins.PostFile(k,v)
	}
	if rep, err = ins.SetTimeout(c.requestTimeout, c.readDataTimeout).Bytes(); err != nil {
	}
	return rep, err
}
