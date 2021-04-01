package main

import (
	"void/voruntime"
)


import "C"

func main(){
	var c chan int=make(chan int,1)
	voruntime.Initrc()
	voruntime.InitUser()
	voruntime.InitInternal()
	voruntime.InitPlugin()
	voruntime.InitSocket()
	<-c
}
