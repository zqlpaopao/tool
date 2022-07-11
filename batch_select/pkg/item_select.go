package pkg

import (
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"strconv"
	"strings"
	"time"
)

//check verification mysql client and order column
func (o *option) check() error {
	if o.mysqlCli == nil {
		return ERRMySqlCli
	}
	if o.orderColumn == "" {
		return ERROrderColumn
	}
	if o.selectFiled != "*" && !strings.Contains(o.selectFiled, o.orderColumn) {
		return ErrSelectNotHaveOrderColumn
	}
	return nil
}

//Run tidy other info
func (o *option) Run() {
	o.getBaseData()
	o.wg.Add(o.handleGoNum)
	for i := 0; i < o.handleGoNum; i++ {
		go o.doing()
	}
	go o.revResCh()
	o.wg.Wait()
	close(o.resCh)
}

//revResCh tidy the result and call back func
func (o *option) revResCh() {
	defer o.savePanicFunc(o.table)
LABEL:
	for {
		select {
		case v, ok := <-o.resCh:
			if !ok {
				break LABEL
			}
			if o.callFunc != nil {
				o.callFunc(v)
				*v = make([]map[string]interface{}, 0, o.limit)
				o.putResMap(*v)
			}
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
	o.revWg.Done()
}

//getBaseData get base data ,will run all
func (o *option) getBaseData() {
	var (
		maxBaseInfo MaxMinInfo
		res         []map[string]interface{}
		err         error
		minMax      *MinMaxInfo
	)
	defer o.savePanicFunc(o.table)
	if maxBaseInfo, err = o.getMaxMinInfo(); nil != err {
		panic(err)
	}
	if maxBaseInfo == (MaxMinInfo{}) {
		panic("get maxMinInfo is fail")
	}
	if res, err = o.getResInfo(o.sqlWhere + " and " + o.orderColumn + " = " + maxBaseInfo.Min); nil != err {
		o.retryFind(nil, err)
		return
	}
	o.resCh <- &res
	o.maxInfo, minMax = maxBaseInfo.Max, o.getMaxId(&MinMaxInfo{}, maxBaseInfo.Min)
	o.minWhereCh <- minMax
}

//running id get max id
func (o *option) getMaxId(info *MinMaxInfo, minId string) *MinMaxInfo {
	var (
		id  int
		err error
	)
	if !o.OrderId {
		info.MinId = minId
		return info
	}
	if id, err = strconv.Atoi(minId); nil != err {
		o.retryFind(nil, err)
	}
	info.MinId, info.MaxId = minId, strconv.Itoa(id+o.limit)
	return info
}

//doing Handle by yourself every goroutine process
func (o *option) doing() {
	defer o.savePanicFunc(o.table)
LABEL:
	for {
		select {
		case v, ok := <-o.minWhereCh:
			if !ok {
				break LABEL
			}
			if v == nil {
				continue
			}
			o.getRes(v)
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
	o.wg.Done()
}

//getRes get every item result
func (o *option) getRes(v *MinMaxInfo) {
	var (
		sqlStr = ""
		err    error
		res    = o.getResMap()
		last   string
	)
	if !o.OrderId {
		sqlStr = o.sqlWhere + " and " + o.orderColumn + " > " + v.MinId
	} else {
		sqlStr = o.sqlWhere + " and " + o.orderColumn + " > " + v.MinId + " and " + o.orderColumn + " <= " + v.MaxId
	}
	if res, err = o.getResInfo(sqlStr); nil != err {
		o.retryFind(v, err)
		return
	}
	if len(res) < 1 && !o.OrderId {
		return
	}
	if len(res) < 1 && o.OrderId {
		o.retryFind(o.getMaxId(v, v.MaxId), nil)
		return
	}
	last = pkg.StringFromAssertionFloat(res[len(res)-1][o.orderColumn])
	o.resCh <- &res
	if last == o.maxInfo {
		close(o.minWhereCh)
		return
	}
	o.retryFind(o.getMaxId(v, last), nil)
}

//retryFind put back the min id
func (o *option) retryFind(id *MinMaxInfo, err error) {
	if id != nil {
		o.minWhereCh <- id
	}
	if nil != err {
		o.err = append(o.err, err)
	}
}
