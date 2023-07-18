package pkg

type Pool[T any] struct {
	Data chan T
}

// GetResMap reusable objects
func (o *Pool[T]) GetResMap() (res T) {
	select {
	case res = <-o.Data:
	// reuse existing buffer
	default:
		// create new buffer
		r := make([]T, 1)
		return r[0]
	}
	return
}

// PutResMap reusable objects
func (o *Pool[T]) PutResMap(res T) {
	select {
	case o.Data <- res:
	// reuse existing buffer
	default:
	}
	return
}
