package pkg

import (
	"errors"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ConnRes interface {
	Close() error
}

//Factory 工厂方法，用于创建连接资源
type Factory func() (ConnRes, error)

type ConnResDoing interface {
	Get() *Conn
	Put(*Conn)
	Release()
	NumPooled() int
}

//Conn 连接资源
type Conn struct {
	idleTimeout time.Time
	conn        ConnRes
}

//连接池
type itemPool struct {
	opt     *option
	clients chan *Conn
	newTag  chan struct{}
}

func NewItemPool(f ...Option) (ConnResDoing, error) {
	i := &itemPool{
		opt:     NewOption(f...),
		clients: nil,
		newTag:  make(chan struct{}),
	}
	i.clients = make(chan *Conn, i.opt.maxNum)
	return i, i.check()
}

//check args
func (i *itemPool) check() error {
	//if i.opt.addr == "" || i.opt.port == "" {
	//	return errors.New("can not addr or port is empty")
	//}
	if i.opt.factory == nil {
		return errors.New("can not conn factory is empty")
	}
	if i.opt.maxNum < 1 {
		return errors.New("can not conn num less 1")
	}
	if i.opt.connTimeout < 1 {
		return errors.New("can not conn time less 1")
	}
	i.makeConnPool()
	return nil
}

//init connPool
func (i *itemPool) makeConnPool() {
	var n int
	for n < i.opt.maxNum {
		if conn := i.makeConn(); conn != nil {
			i.clients <- conn
			n++
		}
	}
}

//new conn
func (i *itemPool) makeConn() (c *Conn) {
	if conn, err := i.opt.factory(); err != nil {
		_ = conn.Close()
		return nil
	} else {
		return &Conn{idleTimeout: time.Now(), conn: conn}
	}
}

//Get  obtain conn
func (i *itemPool) Get() (b *Conn) {
	select {
	case b = <-i.clients:
		if b.idleTimeout.Add(i.opt.connTimeout).Before(time.Now()) {
			return i.Get()
		}
	default:
		b = i.makeConn()
	}
	return

}

//Put return conn
func (i *itemPool) Put(c *Conn) {
	c.Reset()
	select {
	case i.clients <- c:
	default:
		_ = c.conn.Close()
	}
}

//Release resources
//Close connection
func (i *itemPool) Release() {
	close(i.clients)
	for conn := range i.clients {
		_ = conn.conn.Close()
	}
}

//Reset conn
func (c *Conn) Reset() {
	c.idleTimeout = time.Now()
}

//GetHttpConn rev net.Conn
func (c *Conn) GetHttpConn() net.Conn {
	return c.conn.(net.Conn)
}

func (i *itemPool) NumPooled() int {
	return len(i.clients)
}

//GrpcFactor grpc factor
func GrpcFactor(addr string, port string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr+":"+port, grpc.WithInsecure())
}

//HttpFactor http factor
func HttpFactor(tcp, addr string, port string) (net.Conn, error) {
	return net.Dial(tcp, addr+":"+port)
}
