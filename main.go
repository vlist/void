package main

import (
	"void/voruntime"
	"void/voshell"
)

func main(){
	voruntime.Initrc()
	//go voruntime.InitRPC()
	voshell.InitSocket()
}
