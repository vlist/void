package main

import (
	"void/vokernel"
	"void/voruntime"
)

import "C"

func main(){
	var c chan int=make(chan int,1)
	vokernel.InitRC()
	voruntime.InitUserRC()
	voruntime.InitInternal()
	voruntime.Info()
	println()
	voruntime.InitPlugin()
	voruntime.InitSocket()
	<-c
}
