package src

// BufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferPool struct {
	c chan *Task
}

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewBufferPool(size int) (bp *BufferPool) {
	return &BufferPool{
		c: make(chan *Task, size),
	}
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferPool) Get() (b *Task) {
	select {
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = &Task{}
	}
	return
}

// Put returns the given Buffer to the BufferPool.
func (bp *BufferPool) Put(b *Task) {
	b.Reset()
	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}

// NumPooled returns the number of items currently pooled.
func (bp *BufferPool) NumPooled() int {
	return len(bp.c)
}

/*****************************************************ErrorHandler************************************/

// BufferErrPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferErrPool struct {
	c chan *ErrorHandle
}

// NewBufferErrPool creates a new BufferErrPool bounded to the given size.
func NewBufferErrPool(size int) (bp *BufferErrPool) {
	return &BufferErrPool{
		c: make(chan *ErrorHandle, size),
	}
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferErrPool) Get() (b *ErrorHandle) {
	select {
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = &ErrorHandle{}
	}
	return
}

// Put returns the given Buffer to the BufferPool.
func (bp *BufferErrPool) Put(b *ErrorHandle) {
	b.Reset()
	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}

// NumPooled returns the number of items currently pooled.
func (bp *BufferErrPool) NumPooled() int {
	return len(bp.c)
}
