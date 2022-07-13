package pkg

type batchSelect struct {
	batchSelectInfo map[string]*option
}

//NewBatchSelect make new batchSelect
func NewBatchSelect() *batchSelect {
	return &batchSelect{batchSelectInfo: make(map[string]*option)}
}

//InitTask Initialize batch processing hook
func (b *batchSelect) InitTask(task ...InitTaskModel) (err error) {
	for i := 0; i < len(task); i++ {
		if err = task[i].check(); nil != err {
			return
		}
		b.batchSelectInfo[task[i].TaskName] = NewOption(task[i].Opt...)
		if err = b.batchSelectInfo[task[i].TaskName].check(); nil != err {
			return
		}
	}
	return
}

//Run Start multiple consumption tasks by name
func (b *batchSelect) Run(taskName ...string) error {
	for i := 0; i < len(taskName); i++ {
		if _, ok := b.batchSelectInfo[taskName[i]]; !ok {
			return ErrNotHave
		}
		b.batchSelectInfo[taskName[i]].wgAll.Add(1)
		go b.batchSelectInfo[taskName[i]].Run()
	}
	return nil
}

//Wait  the wait all
func (b *batchSelect) Wait() {
	for _, v := range b.batchSelectInfo {
		v.wgAll.Wait()
	}
}
