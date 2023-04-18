package pkg

import (
	"sync/atomic"
)

type BatchHook[T any] struct {
	batchInfo map[string]*OptionItem[T]
}

// NewBatchHook make new BatchHook
func NewBatchHook[T any]() *BatchHook[T] {
	return &BatchHook[T]{batchInfo: make(map[string]*OptionItem[T])}
}

// InitTask Initialize batch processing hook
func (b *BatchHook[T]) InitTask(task InitTaskModel[T]) (err error) {
	if err = task.check(); nil != err {
		return
	}
	b.batchInfo[task.TaskName] = NewOption[T](task.Opt...)
	return nil
}

// Submit Submit multiple tasks by name
func (b *BatchHook[T]) Submit(items SubmitModel[T]) (err error) {
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

// Run Start multiple consumption tasks by name
func (b *BatchHook[T]) Run(taskName string) error {
	if _, ok := b.batchInfo[taskName]; !ok {
		return ErrNotHave
	}
	if b.batchInfo[taskName].IsClose() {
		return ErrTaskClosed
	}
	b.batchInfo[taskName].wg.Add(b.batchInfo[taskName].handleGoNum)
	go b.batchInfo[taskName].Run()
	return nil
}

// Release consumption go pool
func (b *BatchHook[T]) Release(taskName string) error {
	if _, ok := b.batchInfo[taskName]; !ok {
		return ErrNotHave
	}
	if b.batchInfo[taskName].IsClose() {
		return ErrTaskClosed
	}
	atomic.CompareAndSwapInt32(&b.batchInfo[taskName].close, OPENED, CLOSED)
	close(b.batchInfo[taskName].itemCh)
	return nil
}

// WaitAll wait all
func (b *BatchHook[T]) WaitAll() {
	for _, v := range b.batchInfo {
		v.wg.Wait()
	}
}

// Wait wait one
func (b *BatchHook[T]) Wait(taskName string) error {
	if _, ok := b.batchInfo[taskName]; !ok {
		return ErrNotHave
	}
	if b.batchInfo[taskName].IsClose() {
		return ErrTaskClosed
	}
	b.batchInfo[taskName].wg.Wait()
	return nil
}
