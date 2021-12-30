package main

import (
	"github.com/zqlpaopao/tool/downloader/src"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

/**
并发下载大文件
https://mp.weixin.qq.com/s/lUks9X-w0POemthUOORTOQ

*/
//go run . --url https://apache.claz.org/zookeeper/zookeeper-3.7.0/apache-zookeeper-3.7.0-bin.tar.gz
func main() {
	// 默认并发数
	concurrencyN := runtime.NumCPU()
	app := &cli.App{
		Name:  "downloader",
		Usage: "File concurrency downloader",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "`URL` to download",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output `filename`",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"n"},
				Value:   concurrencyN,
				Usage:   "Concurrency `number`",
			},
			&cli.BoolFlag{
				Name:    "resume",
				Aliases: []string{"r"},
				Value:   true,
				Usage:   "Resume download",
			},
		},
		Action: func(c *cli.Context) error {
			strURL := c.String("url")
			filename := c.String("output")
			concurrency := c.Int("concurrency")
			resume := c.Bool("resume")
			return src.NewDownloader(concurrency, resume).Download(strURL, filename)
		},
	}

	//lll
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
