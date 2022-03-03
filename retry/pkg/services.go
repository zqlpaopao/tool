package pkg
//
//import (
//	"time"
//)
//
//type RetryableFunc func() bool
//type onRetryCallbackFun func(uint)
//type onCompleteCallbackFun func(uint, bool, ...interface{})
//
//type services struct {
//	retryConf          *conf                 //重试组件设置
//	onRetryCallback    onRetryCallbackFun    //重试操作触发回调函数
//	onCompleteCallback onCompleteCallbackFun //执行完成触发的回调函数
//}
//
///**
// * 同步方式重试
// * @params retryableFun		需要执行的函数
// * @params CustomParams		自定义参数
// * @return	实际执行次数,最终执行状态
// */
//func (s *services) DoSync(retryableFun RetryableFunc, CustomParams ...interface{}) (uint, bool) {
//	return s.execute(retryableFun, CustomParams...)
//}
//
///**
// * 异步方式重试方法
// * @params retryableFun		需要执行的函数
// * @params CustomParams		自定义参数
// */
//func (s *services) DoAsync(retryableFun RetryableFunc, CustomParams ...interface{}) {
//	go func(retryableFun RetryableFunc, CustomParams ...interface{}) {
//		s.execute(retryableFun, CustomParams...)
//	}(retryableFun, CustomParams...)
//}
//
///**
// * 注册重试回调方法
// * @params retryCallback 重试回调方法
// */
//func (s *services) RegisterRetryCallback(retryCallback onRetryCallbackFun) {
//	s.onRetryCallback = retryCallback
//}
//
///**
// * 注册完成回调方法
// * @params completeCallback 完成回调方法
// */
//func (s *services) RegisterCompleteCallback(completeCallback onCompleteCallbackFun) {
//	s.onCompleteCallback = completeCallback
//}
//
///**
// * 执行重试
// * @params retryableFun 需要执行的函数
// * @params CustomParams	自定义参数
// * @return	实际执行次数,最终执行状态
// */
//func (s *services) execute(retryableFun RetryableFunc, CustomParams ...interface{}) (uint, bool) {
//	var index uint
//	var executeResult = false
//	for index = 1; index <= s.retryConf.attemptsCount; index++ {
//		if s.onRetryCallback != nil {
//			s.onRetryCallback(index)
//		}
//		executeResult = retryableFun()
//		if executeResult {
//			break
//		} else if index == s.retryConf.attemptsCount {
//			break
//		}
//		time.Sleep(s.retryConf.delayType(index, s.retryConf))
//	}
//	if s.onCompleteCallback != nil {
//		s.onCompleteCallback(index, executeResult, CustomParams...)
//	}
//	return index, executeResult
//}
