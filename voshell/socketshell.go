package voshell

import (
	"io"
	"net"
	"os"
	"void/vokernel"
	"void/voruntime"
)

var pa string
func InitSocket(){
	pa=voruntime.RC["socket"]
	println("listening: "+pa)
	os.RemoveAll(pa)
	l,_:=net.Listen("unix",pa)
	for{
		co,_:=l.Accept()
		println("new connection "+pa)
		co.Write([]byte("\r\nconnected to void system socket shell\r\n\r\n"))
		go serve(co)
	}
}
func serve(co net.Conn){
	rline,wline:=io.Pipe()
	var vw=vokernel.VolatileWriter{Destination: wline}
	go func(){
		io.Copy(&vw,co)
		println("disconnected "+pa)
		//if e!=nil{
		//	println(e)
		//}
	}()
	sctx:=vokernel.ShellContext{
		WriterSwitch:     &vw, //writer switch
		Reader: rline,
		Writer: co,
		InternalWriterDestination: wline, //internal receiver
		Delim: '\r',
		Privileged: false,
		Name: pa,
	}
	voruntime.Getsize(sctx)
	for{
		s,e:=sctx.Input(vokernel.Prompt(&sctx))
		if e!=nil{
			println("interrupted")
			sctx.Writer.Close()
			break
		}
		pctx:=voruntime.PreProcess(s,&sctx)
		voruntime.Process(pctx)
	}
}