package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/ssh-cli-pool/pkg"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)


func main() {
	object := "kfj"

	ssh := pkg.NewSsh(
		pkg.WithAddr("11.xx.xx.xx:22"),
		pkg.WithNetwork("tcp"),
		//pkg.WriteTimeout("tcp"),
		pkg.WithSshConfig(&ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password("xxxxx#xxxx"),
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

	for {
		cli, err := p1.Get(object)
		fmt.Println(err)

		fmt.Println(cli.Error())
		fmt.Printf("%#v", cli)

		fmt.Println()
		cli1, err1 := p1.Get(object)
		fmt.Println(err1)

		fmt.Println(cli1.Error())
		fmt.Printf("%#v", cli1)

		time.Sleep(3 * time.Second)
		p1.Put(object, cli)
		p1.Put(object, cli1)

	}

	////go func() {
	////
	////}()
	//
	//for i := 0; i < 3; i++ {
	//	_, _ = sshOld()
	//
	//}
	//for {
	//	fmt.Println(runtime.NumGoroutine())
	//
	//	time.Sleep(5 * time.Second)
	//
	//}

}
