package pkg

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"testing"
	"time"
)

func BenchmarkPool(b *testing.B) {
	t := time.Now()
	object := "kfj"
	var (
		cli  *Ssh
		cli1 *Ssh
		err  error
	)
	ssh := NewSsh(
		WithAddr("11.91.161.27:22"),
		WithNetwork("tcp"),
		//WriteTimeout("tcp"),
		WithSshConfig(&ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password("Meimimazql#3368"),
			},
			HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		}),
		WithCallLog(func(s string) {
			if s != "" {
				fmt.Println("WithCallLog", s)
			}
		}),
	)
	ssh.MakeCli()
	if err = ssh.Error(); err != nil {
		fmt.Println(err)
	}

	p := NewPool[Ssh](
		WithPoolName(object),
		WithItem(ssh),
	)
	p.Submit(ssh)

	p1 := NewObject[*Ssh]().Do()
	if err = p1.Submit((*Pool[*Ssh])(p)); nil != err {
		fmt.Println(err)
	}
	if cli, err = p1.Get(object); nil != err {
		fmt.Println(err)
	}
	//if cli, err = p1.Get(object); nil != err {
	//	fmt.Println(err)
	//}

	p1.Put(object, cli)
	//p1.Put(object, cli1)
	//for {

	for i := 0; i < 100; i++ {
		if cli, err = p1.Get(object); nil != err {
			fmt.Println(err)
		}

		//fmt.Println(cli.Error())
		//fmt.Printf("%#v", cli)

		//fmt.Println()

		if cli1, err = p1.Get(object); err != nil {
			fmt.Println(err)

		}

		//fmt.Println(cli1.Error())
		//fmt.Printf("%#v", cli1)

		//time.Sleep(2 * time.Second)
		p1.Put(object, cli)
		p1.Put(object, cli1)
	}

	fmt.Println(time.Now().Sub(t))
}

func BenchmarkOld(b *testing.B) {
	t := time.Now()
	for i := 0; i < 100; i++ {
		old()
	}
	fmt.Println("old", time.Now().Sub(t))
}

func old() {
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
	_, e := client.NewSession()
	if e != nil {
		fmt.Println(e)
	}

}
