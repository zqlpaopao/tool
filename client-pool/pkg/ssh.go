package pkg

import (
	"golang.org/x/crypto/ssh"
	"net"
	"sync/atomic"
	"time"
)

type Ssh struct {
	client      *ssh.Client
	session     *ssh.Session
	option      *SshOption
	startTime   time.Time
	currentTime time.Time
	err         error
	isUsing     int32
}

func (s *Ssh) IsDelete(maxLifeTime time.Duration) bool {
	return s.startTime.Add(maxLifeTime).After(time.Now())
}

func (s *Ssh) Copy() Item {
	c := *s
	c.option = &SshOption{}
	c.option = s.option
	c.setStartTime()
	c.setCurrentTime()
	atomic.SwapInt32(&s.isUsing, IsUsing)
	c.MakeCli()
	return &c
}

func (s *Ssh) IsUse() bool {
	return atomic.LoadInt32(&s.isUsing) == IsUsing
}

func (s *Ssh) SetUsing() {
	atomic.SwapInt32(&s.isUsing, IsUsing)
}

func (s *Ssh) SetNoUsing() {
	atomic.SwapInt32(&s.isUsing, IsFree)
}

func (s *Ssh) Heartbeat() {
	_, _, err := s.client.Conn.SendRequest("ls -l ", true, []byte("ls -l"))
	if err != nil {
		s.option.callLog(err.Error())
	}
	s.setCurrentTime()
}

func (s *Ssh) setStartTime() {
	s.startTime = time.Now()
}
func (s *Ssh) setCurrentTime() {
	s.currentTime = time.Now()
}

func (s *Ssh) Init() *Ssh {
	s.check()
	s.setStartTime()
	s.MakeCli()
	return s
}

// check Verify necessary parameters
func (s *Ssh) check() {
	if s.option.Addr == "" {
		s.err = ErrAddrEmpty
		return
	}
	if s.option.config == nil {
		s.err = ErrSshClientConfig
		return
	}

}

// MakeCli Create a connection session client
// create a timed out client
func (s *Ssh) MakeCli() {
	if s.err != nil {
		return
	}

	var (
		conn net.Conn
		c    ssh.Conn
		ch   <-chan ssh.NewChannel
		req  <-chan *ssh.Request
	)
	if conn, s.err = net.DialTimeout(
		s.option.network,
		s.option.Addr,
		s.option.dialTimeout); s.err != nil {
		return
	}

	timeoutConn := &Conn{conn, s.option.readTimeout, s.option.writeTimeout}
	if c, ch, req, s.err = ssh.NewClientConn(
		timeoutConn,
		s.option.Addr,
		s.option.config); s.err != nil {
		return
	}

	s.client = ssh.NewClient(c, ch, req)
	s.session, s.err = s.client.NewSession()
}

// Error return err
func (s *Ssh) Error() error {
	return s.err
}
