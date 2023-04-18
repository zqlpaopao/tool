package pkg

const (
	BinBash = "/bin/bash"
	C       = "-c"
	Cmd     = "fping -a -g {ip} -C 1 -i 2 -H 32 -q -t 200 2>&1"

	WorkerNum = 10
	ChanBuf   = 100
	ChanSize  = 10000

	SrcPort = 10000
	DstPort = 6379

	Protocol = 1
	Mark     = "114.114.114.114"
)
