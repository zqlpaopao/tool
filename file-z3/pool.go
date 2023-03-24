package pkg

// GetByteSliceBuf reusable objects
func (o *Option) GetByteSliceBuf() (res []byte) {
	select {
	case res = <-o.byteChan:

	// reuse existing buffer
	default:
		// create new buffer
		return make([]byte, o.byteDataSize)
	}
	return
}

// PutByteSliceBuf reusable objects
func (o *Option) PutByteSliceBuf(res []byte) {
	res = make([]byte, o.byteDataSize)
	select {
	case o.byteChan <- res:
	// reuse existing buffer
	default:
	}
	return
}
