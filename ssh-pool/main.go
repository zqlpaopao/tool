package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/ssh-pool/pkg"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

// Conn wraps a net.Conn, and sets a deadline for every read
// and write operation.
type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}
func SSHDialTimeout(network, addr string, config *ssh.ClientConfig, timeout time.Duration) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println(1)
	timeoutConn := &Conn{conn, timeout, timeout}
	c, chans, reqs, err := ssh.NewClientConn(timeoutConn, addr, config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)

	// this sends keepalive packets every 2 seconds
	// there's no useful response from these, so we can just abort if there's an error
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			fmt.Println("t.C")
			a, b, err := client.Conn.SendRequest("ls -l ", true, []byte("ls -l"))

			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(1, a)
			fmt.Println(2, string(b))
		}
	}()
	s, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		fmt.Println(s)
	}
	return client, nil
}

type SSHKeyboardInteractive map[string]string

func main() {
	t := time.Now()
	object := "kfj"

	ssh := pkg.NewSsh(
		pkg.WithAddr("11.91.161.27:22"),
		pkg.WithNetwork("tcp"),
		//pkg.WriteTimeout("tcp"),
		pkg.WithSshConfig(&ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password("Meimimazql#3368"),
			},
			HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		}),
		pkg.WithCallLog(func(s string) {
			if s != "" {
				fmt.Println("WithCallLog", s)
			}
		}),
	)
	ssh.MakeCli()
	err := ssh.Error()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := pkg.NewPool[pkg.Ssh](
		pkg.WithPoolName(object),
		pkg.WithItem(ssh),
	)
	p.Submit(ssh)

	p1 := pkg.NewObject[*pkg.Ssh]().Do()
	err = p1.Submit((*pkg.Pool[*pkg.Ssh])(p))
	fmt.Println(err)

	fmt.Println("start", time.Now().Sub(t))
	for i := 0; i < 100; i++ {
		cli, _ := p1.Get(object)
		//fmt.Println(err)

		//fmt.Println(cli.Error())
		//fmt.Printf("%#v", cli)

		//fmt.Println()
		cli1, _ := p1.Get(object)
		//fmt.Println(err1)

		//fmt.Println(cli1.Error())
		//fmt.Printf("%#v", cli1)

		p1.Put(object, cli)
		p1.Put(object, cli1)

	}
	fmt.Println("end", time.Now().Sub(t))

	t = time.Now()

	for i := 0; i < 100; i++ {
		sshOld()
	}
	fmt.Println(time.Now().Sub(t))

}

func sshOld() (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("Meimimazql#3368"),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	client, err := ssh.Dial("tcp", "11.91.161.27:22", config)
	if err != nil {
		fmt.Println(err)
	}
	s, e := client.NewSession()
	if e != nil {
		fmt.Println(e)
	}

	//for {
	//	b, er := s.CombinedOutput("ls -l")
	//	fmt.Println(string(b))
	//	fmt.Println(12, er)
	//	s.Stdout = nil
	//	s.Stderr = nil
	//	time.Sleep(3 * time.Second)
	//}

	return s, e
}

func timeout() (*ssh.Client, error) {
	password := "Meimimazql#3368"

	return SSHDialTimeout(
		"tcp",
		"11.91.161.27:22",
		&ssh.ClientConfig{
			Config: ssh.Config{},
			User:   "root",
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			BannerCallback:    nil,
			ClientVersion:     "",
			HostKeyAlgorithms: nil,
			Timeout:           10 * time.Second,
		},
		10*time.Second,
	)

}
