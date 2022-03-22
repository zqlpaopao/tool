package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	config "github.com/zqlpaopao/tool/config/src"
	mysql "github.com/zqlpaopao/tool/mysql/src"
	V2 "gorm.io/gorm"
	"time"
)

const project = "project"


func main(){
	//getWhereSqlByStruct()
	//getWhereSqlBySlice()
	//ReplaceQuestionToDollar()
	GetGormClientV2()

}


// -------------------------------------------------- -------------------------------------------------//
//getWhereSqlByStruct
//注意事项:
//1. between的参数必须两个，即AgeMin/Start和AgeMax/Stop必须同时存在或者同时没有,违反了此条件，where语句依然正确，args会不准确
//2. between参数必须小的声明在大的上面，即不可定义成{AgeMax,AgeMin}
//3. Like参数’*like’匹配%xxx,*like*和’like’匹配%xxx%,’like*’匹配xxx%
//4. 支持指针判定是否存在
// where   addr like ? and desc like ? and job like ? and name = ? and ptr = ? and pt1r = ? or num1 = ? or flo = ?
//[earth% %happ% %engineer ft    0 0]
func getWhereSqlByStruct(){
	type Tmp struct{
		Addr string `column:"and,addr,like*"`
		Desc string `column:"and,desc,like"`
		Job string`column:"and,job,*like"`
		Name string `column:"and,name,="`
		Sal float32 `column:"and,sal,>"`
		AgeMin int`column:"or,age,between"`
		AgeMax int `column:"or,age,between"`
		Start time.Time `column:"and,created,between"`
		Stop time.Time `column:"and,created,between"`
		Jump string `column:"-"`
		Ptr *string `column:"and,ptr,="`
		Ptr1 string `column:"and,pt1r,="`
		Num int64 `column:"or,num,="`
		Num1 *int64 `column:"or,num1,="`
		Float *float32 `column:"or,flo,="`
		Float1 float32 `column:"or,flo1,="`
	}

	te := ""
	var te1 int64 = 0
	var te2 float32 = 0
	var tmp = Tmp{
		Addr:"earth",
		Name:"ft",
		Sal:333,
		AgeMin:9,
		AgeMax:18,
		Desc:"happ",
		Job:"engineer",
		Jump:"jump",
		Ptr: &te,
		Ptr1: " ",
		Num1: &te1,
		Float: &te2,
		Float1: 89,

	}

	fmt.Println(mysql.GenWhereByStruct(tmp,"column"))
}

// getWhereSqlBySlice
//注意事项:如果不添加1=1,如果name为”“切age有值时，就会出现’where or age like ?’
//where   1 = ? and name = ? or name < ? and name > ? and created between ? and ?
func getWhereSqlBySlice(){
	var name = "ft"
	var age =9
	var sal =1000
	var start =  "2018-01-01"
	var stop = "2018-02-02"

	var whereMap = make([][]string,0)
	whereMap = append(whereMap,[]string{
		"","1","=",
	})
	if name!=""{
		whereMap = append(whereMap,[]string{
			"and","name","=",
		})
	}
	if sal !=0 {
		whereMap = append(whereMap,[]string{
			"or","name","<",
		})
	}
	if age!=0{
		whereMap = append(whereMap,[]string{
			"and","name",">",
		})
	}
	if start !="" && stop !=""{
		whereMap =append(whereMap,[]string{
			"and","created","between",
		})
	}
	fmt.Println(mysql.GenWhere(whereMap))
}

// ReplaceQuestionToDollar
//where   1 = $1 and name = $2 or name < $3 and name > $4 and created between $5 and $6
func ReplaceQuestionToDollar() {
	fmt.Println(mysql.ReplaceQuestionToDollar("where   1 = ? and name = ? or name < ? and name > ? and created between ? and ? "))
}

// -------------------------------------------------- -------------------------------------------------//

type CronHost struct {
	Id     int16  `json:"id" cv:"id_1"`
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Port   int32  `json:"port"`
	Remark string `json:"remark"`
}
//Conversion 获取 ON DUPLICATE KEY UPDATE 语句
//获取insert sql 、update sql
func Conversion(){
	c := CronHost{
		Id:     1,
		Name:   "name",
		Alias:  "alias",
		Port:   98,
		Remark: "remark",
	}
	c1,err := mysql.GetSQL(mysql.CoverReqInfo{
		Table: "table",
		StructInfo: &c,
	})
	fmt.Println(err)
	fmt.Println(c1)
}

// -------------------------------------------------- -------------------------------------------------//

//GetCLinet 获取 mysql client
func GetCLinet(){
	var (
		err error
		dbClient *sql.DB
	)
	//初始化环境变量
	if err = config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if dbClient ,err = mysql.Ctx.GetClient(project);nil != err{
		panic(err)
	}

	fmt.Println(dbClient)
}

//GetGormClient 获取 Gorm client
func GetGormClient(){
	var (
		err error
		dbClient *gorm.DB
	)
	//初始化环境变量
	if err = config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if dbClient ,err = mysql.CtxOrm.GetClient(project);nil != err{
		panic(err)
	}

	dbClient = dbClient
}

//GetGormClient 获取 Gorm client
func GetGormClientV2(){
	var (
		err error
		dbClient1 *V2.DB
	)
	//初始化环境变量
	if err = config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if dbClient1 ,err = mysql.CtxOrmV2.GetClient(project);nil != err{
		panic(err)
	}

	err = dbClient1.Debug().Exec("select * from cron_tab").Error
	fmt.Println(err)
	dbClient1 = dbClient1
}
