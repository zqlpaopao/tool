package pkg

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"runtime/debug"
	"testing"
	"time"
)

var DB *gorm.DB

func init() {
	var err error
	dsn := "root:meimima123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type ArpHost struct {
	ArpHostId  uint64    `json:"arp_host_id"`
	Ip         string    `json:"ip"`
	Mac        string    `json:"mac"`
	MacCtime   string    `json:"mac_ctime"`
	MacMtime   string    `json:"mac_mtime"`
	MacTag     string    `json:"mac_tag"`
	IpRole     string    `json:"ip_role"`
	IpCtime    string    `json:"ip_ctime"`
	IpMtime    string    `json:"ip_mtime"`
	DevIp      string    `json:"dev_ip"`
	DevName    string    `json:"dev_name"`
	DevIfShort string    `json:"dev_if_short"`
	DevIfIndex string    `json:"dev_if_index"`
	DevIfName  string    `json:"dev_if_name"`
	Pod        string    `json:"pod"`
	Isp        string    `json:"isp"`
	Idc        string    `json:"idc"`
	Flag       int32     `json:"flag"`
	CreateTime time.Time `json:"create_time"`
}

func(_ ArpHost)TableName()string{
	return "arp_host"
}
func BenchmarkInsert(b *testing.B) {
	t := time.Now()
	for i := 0; i < 100000; i++ {
		item := ArpHost{
			Ip:         "1.1.1.1",
			Mac:        "1111-1111-2222-2222",
			MacCtime:   time.Now().Format("2006-01-02 15:04:05"),
			MacMtime:   time.Now().Format("2006-01-02 15:04:05"),
			MacTag:     "MacTag",
			IpRole:     "IpRole",
			IpCtime:    "IpCtime",
			IpMtime:    "IpMtime",
			DevIp:      "DevIp",
			DevName:    "DevName",
			DevIfShort: "DevIfShort",
			DevIfIndex: "DevIfIndex",
			DevIfName:  "DevIfName",
			Pod:        "Pod",
			Isp:        "Isp",
			Idc:        "Idc",
			Flag:       0,
			CreateTime: time.Now(),
		}

		if err := DB.
			Model(&ArpHost{}).
			//Debug().
			Create(&item).
			Error;nil != err{
				fmt.Println("=========",err)
		}
	}
to := time.Now().Sub(t)
	fmt.Println("insert-HOOK",to)


}

func BenchmarkNewBatchHook(b *testing.B) {
	t := time.Now()
	var err error
	df := NewBatchHook()
	task := []InitTaskModel{
		{
			TaskName: "test1",
			Opt: []Option{
				WithWaitTime(2 * time.Second),
				WithLoopTime(1 * time.Second),
				WithChanSize(1000),
				WithDoingSize(3000),
				WithHandleGoNum(3),
				WithHookFunc(func(item []interface{}) bool {
					//fmt.Println("======len(item)=======", len(item))
					//fmt.Println("=============", i)
					var arpHosts []ArpHost
					for _, vl := range item {
						if v1, ok := vl.(ArpHost); ok {
							arpHosts = append(arpHosts,v1)
						} else {
							os.Exit(1)
						}
					}
					//fmt.Println("len(arpHosts)",len(arpHosts))
					//fmt.Println(len(arpHosts)/3000)
					//os.Exit(2)
					//for i := 0;i < len(arpHosts)/3000;i++{
					//	var info = arpHosts[i*3000:i*3000+3000]
						if err = DB.
							Model(&ArpHost{}).
							//Debug().
							Create(&arpHosts).
							Error;nil != err{
							fmt.Println("=========",err)
							os.Exit(6)
						}

					//}


					return true
				}),
				WithEndHook(func(b bool, i ...interface{}) {
					//fmt.Println("--------------")
					//fmt.Println(b)
					//fmt.Println(len(i))
					//fmt.Println(i)
				}),
				WithSavePanic(func(i interface{}) {
					if err := recover(); err != nil {
						fmt.Println(err)
						fmt.Println(string(debug.Stack()))
						os.Exit(8)
					}
				}),
			},
		},
	}
	
	if err = df.InitTask(task...);nil != err{
		fmt.Println(err)
		os.Exit(3)
	}

	if err := df.Run([]string{"test1"}...); nil != err {
		fmt.Println(err)
		os.Exit(4)
	}

	item := SubmitModel{
		TaskName: "test1",
	}


	for i := 0; i < 100000; i++ {
		items := ArpHost{
			Ip:         "2.2.2.2",
			Mac:        "3333-3333-4444-4444",
			MacCtime:   time.Now().Format("2006-01-02 15:04:05"),
			MacMtime:   time.Now().Format("2006-01-02 15:04:05"),
			MacTag:     "MacTag2",
			IpRole:     "IpRole2",
			IpCtime:    "IpCtime2",
			IpMtime:    "IpMtime2",
			DevIp:      "DevIp2",
			DevName:    "DevName2",
			DevIfShort: "DevIfShort2",
			DevIfIndex: "DevIfIndex2",
			DevIfName:  "DevIfName2",
			Pod:        "Pod2",
			Isp:        "Isp2",
			Idc:        "Idc2",
			Flag:       0,
			CreateTime: time.Now(),
		}
		item.Data = append(item.Data,items)
	}


	if err = df.Submit(item); nil != err {
		fmt.Println(err)
		os.Exit(4)
	}

	if err = df.Release([]string{"test1"}...); nil != err {
		fmt.Println(err)
		os.Exit(5)
	}
	df.WaitAll()
fmt.Println("end")
	to := time.Now().Sub(t)
	fmt.Println("BATCH-HOOK",to)
}


