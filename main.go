package main

import (
	"voidsystem/voruntime"
	"voidsystem/voshell"
)

func main(){
	voruntime.Initrc()
	//go voruntime.InitRPC()
	voshell.InitSocket()
}
