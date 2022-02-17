package src

//Task meta data
type Task struct {
	f   func(string)
	err string
}

//InitTask init task
func InitTask(argF func(err string), err string) *Task {
	return &Task{f: argF, err: err}
}

//Execute run task
func (t *Task) Execute() {
	t.f(t.err)
}

/*****************************************协程池角色*******************************/
type pool struct {
	receiveCh chan *Task
	runCh     chan *Task
	workerNum int
}

//NewPool init pool
func NewPool(n int, size int) *pool {
	return &pool{
		receiveCh: make(chan *Task, size/2-1),
		runCh:     make(chan *Task, size/2-1),
		workerNum: n,
	}
}

//AddTask add task
func (p *pool) AddTask(f *Task) {
	p.receiveCh <- f

}

//worker do task goroutine
func (p *pool) worker(i int) {
	for task := range p.runCh {
		task.Execute()
	}
}

//Run Start the cooperative running task
func (p *pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker(i)
	}
	//Synchronize the received tasks to the internal runChan
	for task := range p.receiveCh {
		p.runCh <- task
	}
}
