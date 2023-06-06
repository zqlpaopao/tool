package pkg

type PoolMan[T any] interface {
	Get() (IdleConn[T], error)
	Put(*IdleConn[T])
	Close(T)
	Release()
	Len() int
}
