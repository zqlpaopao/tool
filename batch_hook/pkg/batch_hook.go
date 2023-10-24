package pkg

type BatchHook[T any] struct {
	taskName string
	batchOpt *OptionItem[T]
}

// NewBatchHook make new BatchHook
func NewBatchHook[T any]() *BatchHook[T] {
	return &BatchHook[T]{}
}

// InitTask Initialize batch processing hook
func (b *BatchHook[T]) InitTask(task InitTaskModel[T]) (err error) {
	if err = task.check(); nil != err {
		return
	}
	b.taskName,
		b.batchOpt =
		task.TaskName,
		NewOption[T](task.Opt...)

	if b.taskName == "" {
		return ERRTaskNameIsEmpty
	}
	return nil
}

// Submit Submit multiple tasks by name
func (b *BatchHook[T]) Submit(items T) {
	b.batchOpt.itemCh <- items
}

// Run Start multiple consumption tasks by name
func (b *BatchHook[T]) Run() {
	b.batchOpt.wg.Add(b.batchOpt.handleGoNum)
	go b.batchOpt.Run()
}

// Release consumption go pool
func (b *BatchHook[T]) Release() {
	close(b.batchOpt.itemCh)
}

// Wait wait all
func (b *BatchHook[T]) Wait() {
	b.batchOpt.wg.Wait()
}
