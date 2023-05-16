package fping

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"
	"sync"
)

type ScanActive struct {
	opt   *Option
	ipCh  chan string
	err   chan error
	errSL []error
	wg    *sync.WaitGroup
}

// NewScanActive init scan object
func NewScanActive(opt *Option) *ScanActive {
	return &ScanActive{opt: opt, ipCh: make(chan string, opt.chanBuf), wg: &sync.WaitGroup{}, err: make(chan error, 10), errSL: make([]error, 0, 20)}
}

// Do doing
func (s *ScanActive) Do() *ScanActive {
	s.wg.Add(s.opt.workerNum)
	for i := 0; i < s.opt.workerNum; i++ {
		go s.loop()
	}
	go s.receiveErr()
	return s
}

// receiveErr receive the errors
func (s *ScanActive) receiveErr() {
	defer s.opt.recover()
LOOP:
	for {
		select {
		case v, ok := <-s.err:
			if !ok {
				break LOOP
			}
			if v == nil {
				continue
			}
			s.errSL = append(s.errSL, v)
		}
	}
}

// loop receive the jobs to worker
func (s *ScanActive) loop() {
	defer s.opt.recover()
LOOP:
	for {
		select {
		case v, ok := <-s.ipCh:
			if !ok {
				break LOOP
			}
			s.scan(v)
		}
	}
	s.wg.Done()
}

// scan scan job
func (s *ScanActive) scan(ipMark string) {
	var (
		cmd     *exec.Cmd
		read    io.ReadCloser
		err     error
		scanner *bufio.Scanner
	)
	cmd = exec.Command(s.opt.binBash, s.opt.C, strings.Replace(s.opt.cmd, "{ip}", ipMark, -1))

	if read, err = cmd.StdoutPipe(); err != nil {
		goto END
	}
	if err = cmd.Start(); nil != err {
		goto END
	}

	scanner = bufio.NewScanner(read)
	for scanner.Scan() {
		s.opt.callback(scanner.Bytes())
	}

	if err = cmd.Wait(); nil != err {
		goto END
	}
	return
END:
	if err == nil {
		return
	}
	s.err <- errors.New(ipMark + "-" + err.Error())
}

func (s *ScanActive) Submit(ip string) {
	s.ipCh <- ip
}

func (s *ScanActive) Release() *ScanActive {
	close(s.ipCh)
	return s
}

func (s *ScanActive) Wait() *ScanActive {
	s.wg.Wait()
	close(s.err)
	return s
}

func (s *ScanActive) Error() []error {
	return s.errSL
}
