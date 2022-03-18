package pkg

import (
	"time"
)

type RetryManager struct {
	retryConf          *option               //重试组件设置
	onRetryCallback    onRetryCallbackFun    //重试操作触发回调函数
	onCompleteCallback onCompleteCallbackFun //执行完成触发的回调函数
}

//NewRetryManager Get retry function entity
func NewRetryManager(opts ...Option) *RetryManager {
	return &RetryManager{retryConf: NewOption(opts...)}
}

//DoSync sync retry
//retryableFun Retry method to execute
//CustomParams Custom parameters
func (r *RetryManager) DoSync(retryableFun RetryableFunc, CustomParams ...interface{}) (uint, bool) {
	return r.execute(retryableFun, CustomParams...)
}

//DoAsync Asynchronous retry method
//retryableFun		Retry method to execute
//CustomParams		Custom parameters
func (r *RetryManager) DoAsync(retryableFun RetryableFunc, CustomParams ...interface{}) {
	go func(retryableFun RetryableFunc, CustomParams ...interface{}) {
		r.execute(retryableFun, CustomParams...)
	}(retryableFun, CustomParams...)
}

//RegisterRetryCallback Register callback method for each retry
//retryCallback  callback method for each retry
func (r *RetryManager) RegisterRetryCallback(retryCallback onRetryCallbackFun) *RetryManager {
	r.onRetryCallback = retryCallback
	return r
}

//RegisterCompleteCallback Registration completion callback method
//CompleteCallback completion callback method
func (r *RetryManager) RegisterCompleteCallback(CompleteCallback onCompleteCallbackFun) *RetryManager {
	r.onCompleteCallback = CompleteCallback
	return r
}

//Execute retry
//@params retryAbleFun function to be executed
//@params CustomParams Custom parameters
//@return Actual execution times, final execution status
func (r *RetryManager) execute(retryableFun RetryableFunc, CustomParams ...interface{}) (uint, bool) {
	var index uint = 1
	var executeResult = false
	for index <=r.retryConf.retryCount{
		if r.onRetryCallback != nil {r.onRetryCallback(index)}
		executeResult = retryableFun()
		if executeResult {break}else if index >= r.retryConf.retryCount{break}
		index++
		time.Sleep(r.retryConf.retryInterval*r.retryConf.delayType(index))
	}
	if r.onCompleteCallback != nil {r.onCompleteCallback(index, executeResult, CustomParams...)}
	return index, executeResult
}
