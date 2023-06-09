package pkg

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Pool Connection pool
type Pool[T any] struct {
	err         error
	Conn        chan *IdleConn[T]
	opt         *Config[T]
	openingConn int64
	close       int32
}

// IdleConn True connection
type IdleConn[T any] struct {
	Conn T
	T    time.Time
}

// InitConn Number of initialized connections
func (p *Pool[T]) InitConn() *Pool[T] {
	for i := int64(0); i < p.opt.InitialCap; i++ {
		var conn T
		if conn, p.err = p.opt.Factory(); nil != p.err {
			continue
		}
		p.Conn <- &IdleConn[T]{Conn: conn, T: time.Now()}
		atomic.AddInt64(&p.openingConn, 1)
	}
	go p.LoopCheck()
	return p
}

// LoopCheck  every check time to check client
func (p *Pool[T]) LoopCheck() {
	if !p.opt.IsCheck {
		return
	}
	var t = time.NewTimer(p.opt.CheckInterval)
LOOP:
	for {
		if atomic.LoadInt32(&p.close) == DefaultIsClose {
			break LOOP
		}
		<-t.C
		p.tidyPool()
		t.Reset(p.opt.CheckInterval)
	}
	_ = t.Stop()
}

// tidyPool check client lifetime
func (p *Pool[T]) tidyPool() {
	var (
		client *IdleConn[T]
		err    error
	)
	if client, err = p.Get(); nil != err {
		goto END
	}
	if client.T.Before(time.Now()) {
		p.Close(client.Conn)
	}
	if err = p.Ping(client.Conn); nil != err {
		p.Close(client.Conn)
		goto END
	}
	p.Put(client)
	return
END:
	if p.opt.Debug && err != nil {
		fmt.Println("tidyPool--->", err)
	}
}

// Get  a connection from the Connection pool
func (p *Pool[T]) Get() (conn *IdleConn[T], err error) {
	var ok bool
	for {
		select {
		case conn, ok = <-p.Conn:
			if !ok {
				err = ErrClosed
				return
			}
			if conn.T.Add(p.opt.IdleTimeout).Before(time.Now()) {
				p.Close(conn.Conn)
				continue
			}
			if err = p.Ping(conn.Conn); err != nil {
				p.Close(conn.Conn)
				continue
			}
			return conn, nil
		default:
			if atomic.LoadInt64(&p.openingConn) > p.opt.MaxCap {
				continue
			}
			cli, er := p.opt.Factory()
			atomic.AddInt64(&p.openingConn, 1)
			return &IdleConn[T]{
				Conn: cli,
				T:    time.Now(),
			}, er
		}
	}
}

// Put  the connection back into the Connection pool
func (p *Pool[T]) Put(conn *IdleConn[T]) {
	if conn == nil {
		return
	}

	if len(p.Conn) > int(p.opt.MaxCap) {
		p.opt.Close(conn.Conn)
		return
	}
	if err := p.Ping(conn.Conn); err != nil {
		p.Close(conn.Conn)
		return
	}
	conn.T = time.Now()
	if p.Conn == nil {
		p.Close(conn.Conn)
		return
	}
	p.Conn <- conn

}

// Close  stop a single connection
func (p *Pool[T]) Close(conn T) {
	atomic.AddInt64(&p.openingConn, -1)
	p.opt.Close(conn)
}

// Ping Check if a single connection is valid
func (p *Pool[T]) Ping(conn T) error {
	return p.opt.Ping(conn)
}

// Release  all connections in the Connection pool
func (p *Pool[T]) Release() {
	close(p.Conn)
	atomic.SwapInt32(&p.close, DefaultIsClose)
	for conn := range p.Conn {
		p.Close(conn.Conn)
	}
}

// Len Existing connections in the Connection pool of len
func (p *Pool[T]) Len() int {
	return len(p.Conn)
}

// Error get error
func (p *Pool[T]) Error() error {
	return p.err
}
