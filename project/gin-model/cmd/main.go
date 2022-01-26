package main

import (
	_ "github.com/zqlpaopao/tool/gin-model/module/init"
	start "github.com/zqlpaopao/tool/gin-model/module/init"
	"sync"
)

func init(){
}


func main(){
	waitGroup := sync.WaitGroup{}
	initWeb(&waitGroup)
	waitGroup.Wait()
}



//initWeb
func initWeb(waitGroup *sync.WaitGroup){
	waitGroup.Add(1)
	go 	start.InitWeb()
}