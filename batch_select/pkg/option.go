package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"runtime/debug"
	"sync"
)

//EndFunc Functions that handle callbacks each time
type EndFunc func(res *[]map[string]interface{})

//SavePanic Functions that handle exception panic
type SavePanic func(i interface{})

type Option interface {
	apply(*option)
}

type option struct {
	debug bool
	//OrderId         bool
	//orderColumnIsTime         bool
	handleGoNum    int
	handleRevGoNum int
	limit          int
	resChanSize    int
	maxInfo        string
	table          string
	orderColumn    string
	sqlWhere       string
	whereCase      []interface{}
	selectFiled    string
	minWhereCh     chan *MinMaxInfo
	resCh          chan *[]map[string]interface{}
	resPool        chan []map[string]interface{}
	err            []error
	callFunc       EndFunc
	savePanicFunc  SavePanic
	wg             *sync.WaitGroup
	revWg          *sync.WaitGroup
	wgAll          *sync.WaitGroup
	mysqlCli       *gorm.DB
}

type OpFunc func(*option)

func NewOption(opt ...Option) *option {
	return clone().WithOptions(opt...)
}

//apply assignment function entity
func (o OpFunc) apply(opt *option) {
	o(opt)
}

//clone  new object
func clone() *option {
	return &option{
		debug:          false,
		limit:          Limit,
		handleGoNum:    HandleGoNum,
		handleRevGoNum: HandleGoNum,
		resChanSize:    ChanSize,
		table:          "",
		orderColumn:    "",
		sqlWhere:       "",
		selectFiled:    "*",
		minWhereCh:     nil,
		resPool:        nil,
		resCh:          nil,
		savePanicFunc:  defaultSavePanic,
		wg:             &sync.WaitGroup{},
		revWg:          &sync.WaitGroup{},
		wgAll:          &sync.WaitGroup{},
		mysqlCli:       nil,
		err:            make([]error, 0, 1),
	}
}

//WithOptions Execute assignment function entity
func (o *option) WithOptions(opt ...Option) *option {
	for _, v := range opt {
		v.apply(o)
	}
	o.initParams()
	return o
}

//initParams Initialization parameters
func (o *option) initParams() {
	o.resCh, o.resPool, o.minWhereCh = make(chan *[]map[string]interface{}, o.resChanSize), make(chan []map[string]interface{}, o.handleGoNum/3), make(chan *MinMaxInfo, o.handleGoNum)
	if o.debug {
		o.mysqlCli = o.mysqlCli.Debug()
	}
}

//WithDebug debug mysql
func WithDebug(debug bool) OpFunc {
	return func(o *option) {
		o.debug = debug
	}
}

//WithLimit How much to start processing default 100
func WithLimit(size int) OpFunc {
	return func(o *option) {
		o.limit = size
	}
}

//WithHandleGoNum Number of goroutine processed default 100
func WithHandleGoNum(num int) OpFunc {
	return func(o *option) {
		o.handleGoNum = num
	}
}

//WithHandleRevGoNum Number of goroutine rev
func WithHandleRevGoNum(num int) OpFunc {
	return func(o *option) {
		o.handleRevGoNum = num
	}
}

//WithTable will doing table name
func WithTable(tableName string) OpFunc {
	return func(o *option) {
		o.table = tableName
	}
}

//WithOrderColumn order column
func WithOrderColumn(name string) OpFunc {
	return func(o *option) {
		o.orderColumn = name
	}
}

//WithSqlWhere where
func WithSqlWhere(sqlWhere string, whereCase []interface{}) OpFunc {
	return func(o *option) {
		o.sqlWhere, o.whereCase = sqlWhere, whereCase
	}
}

//WithSelectFiled select filed
func WithSelectFiled(filed string) OpFunc {
	return func(o *option) {
		o.selectFiled = filed
	}
}

//WithMysqlSqlCli where
func WithMysqlSqlCli(cli *gorm.DB) OpFunc {
	return func(o *option) {
		o.mysqlCli = cli
	}
}

//WithResChanSize rev chan size
func WithResChanSize(size int) OpFunc {
	return func(o *option) {
		o.resChanSize = size
	}
}

//WithCallFunc callback end func
func WithCallFunc(endHook EndFunc) OpFunc {
	return func(o *option) {
		o.callFunc = endHook
	}
}

//defaultSavePanic
func defaultSavePanic(i interface{}) {
	if err := recover(); nil != err {
		fmt.Println(i)
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
	}
}
