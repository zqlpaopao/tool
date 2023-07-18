package etcd

import (
	"context"
	"errors"
	"fmt"
	format "github.com/zqlpaopao/tool/format/src"
	"github.com/zqlpaopao/tool/register-discovery/pkg"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"sync/atomic"
	"time"
)

type DoEtcd struct {
	err     error
	etcdCLi *etcd.Client
	option  *Options
	close   atomic.Bool
}

func NewPDoEtcd(etcdCLi *etcd.Client,
	opt ...OptionFunc) *DoEtcd {
	o := &DoEtcd{etcdCLi: etcdCLi,
		option: NewOptions(opt...),
		close:  atomic.Bool{}}
	return o
}

func (e *DoEtcd) Register(info *pkg.RegisterInfo) {
	var (
		kv            etcd.KV
		lease         etcd.Lease
		leaseGrantRes *etcd.LeaseGrantResponse
		leaseId       etcd.LeaseID
		keepResChan   <-chan *etcd.LeaseKeepAliveResponse
		keepRes       *etcd.LeaseKeepAliveResponse
		ctx           context.Context
		cancelFunc    context.CancelFunc
	)
	// make lease
	lease = etcd.NewLease(e.etcdCLi)
	//申请5s的租约
	if leaseGrantRes, e.err = lease.Grant(context.TODO(), e.option.Registerer.PushTime); nil != e.err {
		e.option.CallBackErr("Register.lease.Grant", e.err)
		return
	}

	//拿到租约id
	leaseId = leaseGrantRes.ID

	//准备取消续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//确保函数推出后，自动停止续租
	defer cancelFunc()

	//立即释放租约
	defer func() {
		if _, e.err = lease.Revoke(context.TODO(), leaseId); nil != e.err {
			e.option.CallBackErr("Register.defer.lease.Revoke", e.err)
			return
		}
	}()

	//自动续租
	if keepResChan, e.err = lease.KeepAlive(ctx, leaseId); nil != e.err {
		e.option.CallBackErr("Register.lease.KeepAlive", e.err)
		return
	}

	//处理应答续租的协程
	go func() {
		for {
			if e.close.Load() {
				goto END
			}
			select {
			case keepRes = <-keepResChan:
				if keepResChan == nil {
					e.option.CallBackErr("Register.keepResChan.keepResChan", errors.New("keepResChan is nil,lease is over"))
					goto END
				}
				e.Debug(keepRes.String())
			}
		}
	END:
	}()

	//如果锁设置不成功，then 设置它，else 设置失败
	kv = etcd.NewKV(e.etcdCLi)

	if _, e.err = kv.Put(context.TODO(), e.option.Registerer.Addr, info.RegisterInfo, etcd.WithLease(leaseId)); nil != e.err {
		e.option.CallBackErr("Register.kv.Put", e.err)
		return
	}

	for {
		if e.close.Load() {
			goto ENDS
		}
		time.Sleep(e.option.PullTime)
	}
ENDS:
}

func (e *DoEtcd) UnRegister() {
	e.close.Store(true)
}

func (e *DoEtcd) Discovery() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		resp   *etcd.GetResponse
	)
	ctx, cancel = context.WithTimeout(context.Background(), e.option.Discovery.PullTime)
	defer cancel()

	if resp, e.err = e.etcdCLi.Get(ctx, e.option.Discovery.Addr); nil != e.err {
		e.option.CallBackErr("Discovery.etcdCLi.Get", e.err)
		return
	}
	e.option.CallBack(resp)
	go e.Watch()
}

func (e *DoEtcd) Watch() {
	var (
		resp  etcd.WatchChan
		watch etcd.Watcher
	)
	watch = etcd.NewWatcher(e.etcdCLi)
	go e.UnWatch(watch)
	resp = watch.Watch(context.Background(), e.option.Discovery.Addr) // <-chan WatchResponse
	for res := range resp {
		for _, ev := range res.Events {
			//fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			e.option.CallBackWatch(ev)
		}
	}
}

func (e *DoEtcd) UnWatch(watch etcd.Watcher) {
LOOP:
	for {
		if e.close.Load() {
			break LOOP
		}
	}
	if err := watch.Close(); nil != err {
		e.option.CallBackErr("UnWatch ", err)
	}

}

func (e *DoEtcd) UnDiscovery() {
	e.close.Store(true)
}

func (e *DoEtcd) Error() error {
	return e.err
}

func (e *DoEtcd) Debug(info string) {
	if e.option.Debug {
		return
	}
	log.Println("DEBUG " + info)
}

func CallBackErr(funcName string, err error) {
	if err != nil {
		format.PrintRed(funcName + err.Error())
	}
}

func CallBackWatch(info *etcd.Event) {
	fmt.Printf("%#v", info)
}

func CallBack(info *etcd.GetResponse) {
	fmt.Printf("%#v", info)
}
