package pkg

type BatchSelect[T any] struct {
	taskName string
	batchOpt *Option[T]
}

// NewBatchSelect make new BatchSelect
func NewBatchSelect[T any]() *BatchSelect[T] {
	return &BatchSelect[T]{}
}

// InitTask Initialize batch processing hook
func (b *BatchSelect[T]) InitTask(task InitTaskModel[T]) (err error) {
	if err = task.check(); nil != err {
		return
	}

	b.taskName,
		b.batchOpt =
		task.TaskName,
		NewOption[T](task.Opt...)

	return b.batchOpt.check()
}

// Run Start multiple consumption tasks by name
func (b *BatchSelect[T]) Run(taskName string) error {
	if taskName != b.taskName {
		return ErrNotHave
	}

	b.batchOpt.wgAll.Add(1)
	go b.batchOpt.Run()
	return nil
}

// Wait  the wait all
func (b *BatchSelect[T]) Wait() {
	b.batchOpt.wgAll.Wait()
}
