package pkg

// Get reusable objects
func (o *Option[T]) getResMap() (res []T) {
	select {
	case res = <-o.resPool:
	// reuse existing buffer
	default:
		// create new buffer
		return make([]T, 0, o.limit)
	}
	return
}

// Put reusable objects
func (o *Option[T]) putResMap(res []T) {
	select {
	case o.resPool <- res:
	// reuse existing buffer
	default:
	}
	return
}
