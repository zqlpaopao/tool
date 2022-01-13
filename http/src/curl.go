package src

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type curlManager struct {
	request *http.Request
	client  http.Client
}

var (
	CurlCtx = new(curlManager)
)

//NewRequest get request client
func (c *curlManager) NewRequest(args RequestArgs) (curlMan *curlManager, err error) {
	if err = chekArgs(&args); nil != err {
		return
	}
	curlMan = &curlManager{
		request: &http.Request{Header: http.Header{}},
		client: http.Client{},
	}
	data := url.Values{}
	curlMan.assignReqParam(&data, args.RequestParams)
	curlMan.SetTimeOut(args.SetTimeOut)
	err = curlMan.switchMethod(&args,&data)
	curlMan.client = http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	return curlMan, err
}

func(c *curlManager)switchMethod(args *RequestArgs,data *url.Values)(err error){
	var body io.Reader
	switch args.RequestType{
	case RequestPost:
		c.assignHeader(args.RequestHeader)
		body = strings.NewReader(data.Encode())
	case RequestGet:
		args.RequestUrl = tidyGetArgs(args.RequestUrl,data.Encode())
		body = nil
	}
	c.request, err = http.NewRequest(args.RequestType.String(), args.RequestUrl, body)
	return
}


func tidyGetArgs(url string,args string)string{
	if strings.Contains(url,"?"){
		return url+"&"+args
	}
	return url+"?"+args
}

//chekArgs args check
func chekArgs(args *RequestArgs) error {
	if args.RequestUrl == "" {
		return errors.New("request url is empty")
	}
	if args.RequestType < 0 {
		return errors.New("request mode error")
	}
	if args.RequestHeader == nil {
		args.RequestHeader = map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}
	}
	return nil
}

//assignReqParam assignment request parameters
func (c *curlManager) assignReqParam(data *url.Values, args map[string]string) {
	for k, v := range args {
		data.Set(k, v)
	}
}

//assignHeader assignment request header parameters
func (c *curlManager) assignHeader(args map[string]string) {
	for k, v := range args {
		c.request.Header.Set(k, v)
	}
}

//SetTimeOut Set timeout
func (c *curlManager) SetTimeOut(timeOut int64) {
	if timeOut > 0 {
		c.client.Timeout = time.Duration(timeOut) * time.Second
	} else {
		c.client.Timeout = time.Duration(10) * time.Second
	}
}

//SetHeader Add request header
func (c *curlManager) SetHeader(key, value string) {
	c.request.Header.Add(key, value)
}

//SetContentType set ContentType
func (c *curlManager) SetContentType(contentType string) {
	c.request.Header.Set("Content-Type", contentType)
}

//Send request
func (c *curlManager) Send() (responseData []byte, err error) {
	responseRe, err := c.client.Do(c.request)
	if err != nil {
		return
	}
	defer func(responseRe *http.Response) {
		_ = responseRe.Body.Close()
	}(responseRe)
	responseData, err = ioutil.ReadAll(responseRe.Body)
	if err != nil {
		return
	}
	if responseRe.StatusCode != http.StatusOK {
		err = fmt.Errorf("RequestCodeError:%v", responseRe.Status)
		return
	}
	return
}
