package pkg

import (
	"sync/atomic"
)

type BatchHook struct {
	batchInfo map[string]*option
}

//NewBatchHook make new BatchHook
func NewBatchHook() *BatchHook {
	return &BatchHook{batchInfo: make(map[string]*option)}
}

//InitTask Initialize batch processing hook
func (b *BatchHook) InitTask(task ...InitTaskModel) (err error) {
	for i := 0; i < len(task); i++ {
		if err = task[i].check(); nil != err {
			return
		}
		b.batchInfo[task[i].TaskName] = NewOption(task[i].Opt...)
	}
	return nil
}

//Submit Submit multiple tasks by name
func (b *BatchHook) Submit(items SubmitModel) (err error) {
	if _, ok := b.batchInfo[items.TaskName]; !ok {
		return ErrNotHave
	}
	if b.batchInfo[items.TaskName].IsClose() {
		return ErrTaskClosed
	}
	if err = b.batchInfo[items.TaskName].check(); nil != err {
		return
	}
	for i := 0; i < len(items.Data); i++ {
		b.batchInfo[items.TaskName].itemCh <- items.Data[i]
	}
	return
}

//Run Start multiple consumption tasks by name
func (b *BatchHook) Run(taskName ...string) error {
	for i := 0; i < len(taskName); i++ {
		if _, ok := b.batchInfo[taskName[i]]; !ok {
			return ErrNotHave
		}
		if b.batchInfo[taskName[i]].IsClose() {
			return ErrTaskClosed
		}
		b.batchInfo[taskName[i]].wg.Add(b.batchInfo[taskName[i]].handleGoNum)
		go b.batchInfo[taskName[i]].Run()
	}
	return nil
}

//Release consumption go pool
func (b *BatchHook) Release(taskName ...string) error {
	for i := 0; i < len(taskName); i++ {
		if _, ok := b.batchInfo[taskName[i]]; !ok {
			return ErrNotHave
		}
		if b.batchInfo[taskName[i]].IsClose() {
			return ErrTaskClosed
		}
		atomic.CompareAndSwapInt32(&b.batchInfo[taskName[i]].close, OPENED, CLOSED)
		close(b.batchInfo[taskName[i]].itemCh)
	}
	return nil
}

//WaitAll wait all
func (b *BatchHook) WaitAll() {
	for _, v := range b.batchInfo {
		v.wg.Wait()
	}
}

//Wait wait one
func (b *BatchHook) Wait(taskName ...string) error {
	for i := 0; i < len(taskName); i++ {
		if _, ok := b.batchInfo[taskName[i]]; !ok {
			return ErrNotHave
		}
		if b.batchInfo[taskName[i]].IsClose() {
			return ErrTaskClosed
		}
		b.batchInfo[taskName[i]].wg.Wait()
	}
	return nil
}
