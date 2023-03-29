package pkg

// GetByteSliceBuf reusable objects
func (o *Option) GetByteSliceBuf() []byte {
	select {
	case res := <-o.byteChan:
		return *res
	// reuse existing buffer
	default:
		// create new buffer
		return make([]byte, o.ByteDataSize)
	}
}

// PutByteSliceBuf reusable objects
func (o *Option) PutByteSliceBuf(res *[]byte) {
	*res = make([]byte, o.ByteDataSize)
	select {
	case o.byteChan <- res:
	// reuse existing buffer
	default:

	}
	return
}
