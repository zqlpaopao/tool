package pkg

// Pool is a generic data buffer pool that initializes specific types and data before use
// Data is the type of cached data
// New that needs to be created when the cache does not exist
type Pool[T any] struct {
	Data chan *T
	New  func() *T
}

// Get  the data
func (p *Pool[T]) Get() (res *T) {
	select {
	case res, _ = <-p.Data:
		return
	default:
		return p.New()
	}
}

// Put the data to cached
func (p *Pool[T]) Put(res *T) {
	select {
	case p.Data <- res:
		return
	default:

	}
}
