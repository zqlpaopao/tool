package pkg

//Get reusable objects
func (o *option) getResMap() (res []map[string]interface{}) {
	select {
	case res = <-o.resPool:
	// reuse existing buffer
	default:
		// create new buffer
		return make([]map[string]interface{}, 0, o.limit)
	}
	return
}

//Put reusable objects
func (o *option) putResMap(res []map[string]interface{}) {
	select {
	case o.resPool <- res:
	// reuse existing buffer
	default:
	}
	return
}
