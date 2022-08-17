package pkg

import (
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
	go o.producer()
	o.consumer()
}

//getBaseData get base data ,will run all
func (o *option) getBaseData() {
	var (
		maxBaseInfo MaxMinInfo
		err         error
	)
	defer o.savePanicFunc(o.table)
	if maxBaseInfo, err = o.getMaxMinInfo(); nil != err {
		panic(err)
	}
	if maxBaseInfo == (MaxMinInfo{}) {
		close(o.minWhereCh)
		return
	}
	if err = o.OrderIdDebug(&maxBaseInfo); nil != err {
		panic(err)
	}
	return
}

//OrderIdDebug If it is obtained according to the self increasing ID, create an acquisition interval
func (o *option) OrderIdDebug(maxBaseInfo *MaxMinInfo) (err error) {
	var minId, maxId, loop int
	if minId, err = strconv.Atoi(maxBaseInfo.Min); nil != err {
		return
	}
	if maxId, err = strconv.Atoi(maxBaseInfo.Max); nil != err {
		return
	}
	loop = (maxId - minId + o.limit) / o.limit
	go o.sendMinWhere(minId, loop)
	return
}

//sendMinWhere Send to producer
func (o *option) sendMinWhere(minId, loop int) {
	for i := 1; i <= loop; i++ {
		info := &MinMaxInfo{
			MinId: strconv.Itoa(minId),
			MaxId: strconv.Itoa(minId + o.limit),
		}
		o.minWhereCh <- info
		minId = minId + o.limit
	}
	close(o.minWhereCh)
}

//doing Handle by yourself every goroutine process
func (o *option) producer() {
	o.wg.Add(o.handleGoNum)
	for i := 0; i < o.handleGoNum; i++ {
		go o.producerLoopWorker()
	}
	o.wg.Wait()
	close(o.resCh)
}

//producerLoopWorker Producers obtain production interval data
func (o *option) producerLoopWorker() {
	defer o.savePanicFunc("producerLoopWorker-" + o.table)
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
			o.producerDbData(v)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	o.wg.Done()
}

//getRes get every item result
func (o *option) producerDbData(v *MinMaxInfo) {
	var (
		sqlStr = ""
		err    error
		res    = o.getResMap()
	)
	sqlStr = o.orderColumn + " >= '" + v.MinId + "' and " + o.orderColumn + " < " + "'" + v.MaxId + "'"
	if "" != strings.Trim(o.sqlWhere, " ") {
		sqlStr += " and " + o.sqlWhere
	}
	if res, err = o.getResInfo(sqlStr); nil != err {
		o.retryFind(v, err)
		return
	}
	if len(res) < 1 {
		return
	}
	o.resCh <- &res
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

//consumer
func (o *option) consumer() {
	o.revWg.Add(o.handleRevGoNum)
	for i := 0; i < o.handleRevGoNum; i++ {
		go o.revResCh()
	}
	o.revWg.Wait()
	o.wgAll.Done()
}

//revResCh tidy the result and call back func
func (o *option) revResCh() {
	defer o.savePanicFunc("revResCh-" + o.table)
LABEL:
	for {
		select {
		case v, ok := <-o.resCh:
			if !ok {
				break LABEL
			}
			if o.callFunc != nil && v != nil && len(*v) > 0 {
				o.callFunc(v)
				*v = make([]map[string]interface{}, 0, o.limit+1)
				o.putResMap(*v)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	o.revWg.Done()
}
