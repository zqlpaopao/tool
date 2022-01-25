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

//00cfa701-8ee8-02a6-8494-89a478f02cb1 172.21.177.93
//b5a2e387-b83e-93ba-9710-c81014c74bfd 172.21.177.94
type DevIfCmdb struct {
	DevIfId    string    `json:"dev_if_id"`
	DevIp      string    `json:"dev_ip"`
	DeviceId   string    `json:"device_id"`
	DevIfIndex string    `json:"dev_if_index"`
	DevIfAlias string    `json:"dev_if_alias"`
	DevIfDesc  string    `json:"dev_if_desc"`
	DevIfName  string    `json:"dev_if_name"`
	Tag        string    `json:"tag"`
	CTime     string `json:"c_time"`
	UTime      string `json:"u_time"`
	CreateTime string `json:"create_time"`
}

func(_ DevIfCmdb)TableName()string{
	return "dev_if_cmdb"
}


func GetPorts1(devId string,ip string){
	url := "http://cmdb.jd.com/v1.1/network_devices/ports?user=cmdb_all&timestamp=1&auth=1"
	method := "POST"


	ip = "172.21.177.10"

	payload := strings.NewReader(`
{"device_ids":["4f4ae49b-b535-d325-23fe-c7e4c32045cd"]}`)

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


	if resp.Code != 2000{
		return
	}

	if len(resp.Data.List) < 1{
		fmt.Println("获取为空的id",devId)
		return
	}

	var (
		db  *gorm.DB
	)



	if db, err = src2.NewSSHGormClient(&src2.Config{
		Addr:   "11.91.161.27:22",
		User:   "root",
		Passwd: "@Noah0b2",
	}, &src2.MysqlConfig{
		UserName: "joypaw",
		PassWd:   "^1Joy_paw9$",
		IpPort:   "172.18.145.107:3306",
		Dbname:   "joypaw_base",
	}); nil != err {
		panic(err)
	}





	db = db.Debug()
	err = db.Exec("delete from dev_if_cmdb where dev_ip = ?",ip).Error
	fmt.Println(err)

	//src2 "github.com/zqlpaopao/tool/ssh-tool/src"

	for _ ,v := range resp.Data.List{
		 m := DevIfCmdb{
			 DevIfId:    v.Id,
			 DevIp:      ip,
			 DeviceId:   v.DeviceId,
			 DevIfIndex: v.SnmpIfindex,
			 DevIfAlias: v.Alias,
			 DevIfDesc:  v.Name,
			 DevIfName:  v.Name,
			 Tag:        "",
			 CTime:      time.Now().Format("2006-01-02 15:04:05"),
			 UTime:      time.Now().Format("2006-01-02 15:04:05"),
			 CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		 }
		 if err = db.Create(&m).Error;err != nil{
			 fmt.Println(err)
			 continue
		 }
	}





	return
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


