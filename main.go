package main

import (
	"void/voruntime"
)

func main(){
	var c chan int=make(chan int,1)
	voruntime.Initrc()
	voruntime.InitInternal()
	voruntime.InitSocket()
	<-c
}
