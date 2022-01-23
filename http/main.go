package main

import (
	"encoding/json"
	"fmt"
	src2 "github.com/zqlpaopao/tool/ssh-tool/src"

	//"git.jd.com/npd_automation/joypaw-link-detection/common"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/jinzhu/gorm"
	"github.com/zqlpaopao/tool/http/src"

	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//https://beego.vip/docs/module/httplib.md
func main() {
	//GetPorts1("","172.21.177.94")
	YuanSheng()
}

// YuanSheng post 自己拼接参数[]byte类型
//get 可传入
func YuanSheng() {
	url := "http://xx.xx.xx.com"
	b , err := json.Marshal(&map[string][]string{
		"device_ids":[]string{
			"4f4ae49b-b535-d325-23fe-c7e4c32045cd",
		},
	})


	c, err := src.CurlCtx.NewRequest(src.RequestArgs{
		RequestUrl: url,
		RequestParamsGet: map[string]string{
			"device_ids" :"v",
		},
		RequestParamsPost: b,
		RequestType:   src.RequestPost,
		RequestHeader: src.RequestRaw.String(),
		SetTimeOut:    10,
	})

	fmt.Printf("%#v",c)

	fmt.Println(err)

	r, err := c.Send()
	fmt.Println(string(r))
	fmt.Println(err)
}
































type List struct {
	PortType        string `json:"port_type"`
	PortRole        string `json:"port_role"`
	DeviceId        string `json:"device_id"`
	DeviceType      string `json:"device_type"`
	Name            string `json:"name"`
	Alias           string `json:"alias"`
	Describe        string `json:"describe"`
	AdminStatus     int    `json:"admin_status"`
	OperationStatus int    `json:"operation_status"`
	SnmpIfindex     string `json:"snmp_ifindex"`
	Mode            string `json:"mode"`
	Speed           int    `json:"speed"`
	CreatedTime     string `json:"created_time"`
	UpdatedTime     string `json:"updated_time"`
	Aid             int    `json:"aid"`
	Id              string `json:"id"`
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TotalCount int `json:"total_count"`
		Page       int `json:"page"`
		Length     int `json:"length"`
		List       []List
	} `json:"data"`
}
type T struct {
	DeviceIds []string `json:"device_ids"`
}

func GetPorts(devId string){
	url := "http://cmdb.jd.com/v1.1/network_devices/ports?user=cmdb_all&timestamp=1&auth=1"
	method := "POST"

	var m T
	m.DeviceIds = []string{
		devId,
	}

	b , err := json.Marshal(m)
	fmt.Println(string(b))

	payload := strings.NewReader(string(b))

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	 var resp Resp
	if err := json.Unmarshal(body,&resp);err != nil{
		fmt.Println(err)
	}

	fmt.Println(len(resp.Data.List))

}
func libHttp() {
	c, err := httplib.Get("https://baidu.com").Debug(true).Response()
	fmt.Println(c)
	fmt.Println(err)
}


