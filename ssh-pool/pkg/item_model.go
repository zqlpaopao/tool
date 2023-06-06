package pkg

import (
	"time"
)

// Item To implement this interface, you can use pool
type Item interface {
	//Copy When there are no twin objects connected
	Copy() Item
	//MakeCli Create your own suitable connection
	MakeCli()
	//IsUse Is it in use
	IsUse() bool
	//SetUsing set it in use
	SetUsing()
	//SetNoUsing set it not use
	SetNoUsing()
	//IsDelete Do you need to delete it
	IsDelete(duration time.Duration) bool
	//Heartbeat heart beat detection
	Heartbeat()
}
