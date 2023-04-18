package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"runtime/debug"
	"sync"
)

// EndFunc Functions that handle callbacks each time
type EndFunc[T any] func(res *[]T)

// SavePanic Functions that handle exception panic
type SavePanic func(i ...interface{})

type OptionInter[T any] interface {
	apply(*Option[T])
}

type Option[T any] struct {
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
	resCh          chan *[]T
	resPool        chan []T
	err            []error
	callFunc       EndFunc[T]
	savePanicFunc  SavePanic
	wg             *sync.WaitGroup
	revWg          *sync.WaitGroup
	wgAll          *sync.WaitGroup
	mysqlCli       *gorm.DB
}

type OpFunc[T any] func(*Option[T])

func NewOption[T any](opt ...OptionInter[T]) *Option[T] {
	return clone[T]().WithOptions(opt...)
}

// apply assignment function entity
func (o OpFunc[T]) apply(opt *Option[T]) {
	o(opt)
}

// clone  new object
func clone[T any]() *Option[T] {
	return &Option[T]{
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

// WithOptions Execute assignment function entity
func (o *Option[T]) WithOptions(opt ...OptionInter[T]) *Option[T] {
	for _, v := range opt {
		v.apply(o)
	}
	o.initParams()
	return o
}

// initParams Initialization parameters
func (o *Option[T]) initParams() {
	o.resCh, o.resPool, o.minWhereCh = make(chan *[]T, o.resChanSize), make(chan []T, o.handleGoNum/3), make(chan *MinMaxInfo, o.handleGoNum)
	if o.debug {
		o.mysqlCli = o.mysqlCli.Debug()
	}
}

// WithDebug debug mysql
func WithDebug[T any](debug bool) OpFunc[T] {
	return func(o *Option[T]) {
		o.debug = debug
	}
}

// WithLimit How much to start processing default 100
func WithLimit[T any](size int) OpFunc[T] {
	return func(o *Option[T]) {
		o.limit = size
	}
}

// WithHandleGoNum Number of goroutine processed default 100
func WithHandleGoNum[T any](num int) OpFunc[T] {
	return func(o *Option[T]) {
		o.handleGoNum = num
	}
}

// WithHandleRevGoNum Number of goroutine rev
func WithHandleRevGoNum[T any](num int) OpFunc[T] {
	return func(o *Option[T]) {
		o.handleRevGoNum = num
	}
}

// WithTable will doing table name
func WithTable[T any](tableName string) OpFunc[T] {
	return func(o *Option[T]) {
		o.table = tableName
	}
}

// WithOrderColumn order column
func WithOrderColumn[T any](name string) OpFunc[T] {
	return func(o *Option[T]) {
		o.orderColumn = name
	}
}

// WithSqlWhere where
func WithSqlWhere[T any](sqlWhere string, whereCase []interface{}) OpFunc[T] {
	return func(o *Option[T]) {
		o.sqlWhere, o.whereCase = sqlWhere, whereCase
	}
}

// WithSelectFiled select filed
func WithSelectFiled[T any](filed string) OpFunc[T] {
	return func(o *Option[T]) {
		o.selectFiled = filed
	}
}

// WithMysqlSqlCli where
func WithMysqlSqlCli[T any](cli *gorm.DB) OpFunc[T] {
	return func(o *Option[T]) {
		o.mysqlCli = cli
	}
}

// WithResChanSize rev chan size
func WithResChanSize[T any](size int) OpFunc[T] {
	return func(o *Option[T]) {
		o.resChanSize = size
	}
}

// WithCallFunc callback end func
func WithCallFunc[T any](endHook EndFunc[T]) OpFunc[T] {
	return func(o *Option[T]) {
		o.callFunc = endHook
	}
}

// defaultSavePanic
func defaultSavePanic(i ...interface{}) {
	if err := recover(); nil != err {
		fmt.Println(i)
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
	}
}
