package src

// Task meta data
type Task struct {
	f   func(string)
	err string
}

// resource pool Reuse objects and avoid frequent affectionate objects
var buffSize *BufferPool

// NewTaskBufferPool creates a new BufferPool bounded to the given size.
func NewTaskBufferPool(i int) {
	buffSize = NewBufferPool(i)
}

// Reset reset resource
func (t *Task) Reset() {
	t.f, t.err = nil, ""
}

// InitTask init task
func InitTask(argF func(err string), err string) (task *Task) {
	task = buffSize.Get()
	task.f, task.err = argF, err
	return task
}

// Execute run task
func (t *Task) Execute() {
	t.f(t.err)
	buffSize.Put(t)
}

// Pool *****************************************协程池角色*******************************/
type Pool struct {
	receiveCh chan *Task
	runCh     chan *Task
	workerNum int
}

// NewPool init Pool
func NewPool(n int, size int) *Pool {
	return &Pool{
		receiveCh: make(chan *Task, size/2-1),
		runCh:     make(chan *Task, size/2-1),
		workerNum: n,
	}
}

// AddTask add task
func (p *Pool) AddTask(f *Task) {
	p.receiveCh <- f

}

// worker do task goroutine
func (p *Pool) worker(_ int) {
	for task := range p.runCh {
		task.Execute()
	}
}

// Run Start the cooperative running task
func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker(i)
	}
	//Synchronize the received tasks to the internal runChan
	for task := range p.receiveCh {
		p.runCh <- task
	}
}
