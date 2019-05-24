package main

import (
	"adapter/modules/logger"
	"time"
)

func main(){
	for i:=1; i<=10; i++ {
		logger.Sendlog("test.log", "info", "hello word!", i)
		time.Sleep(time.Second*1)
	}
}
