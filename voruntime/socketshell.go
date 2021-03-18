package voruntime

import (
	"io"
	"log"
	"net"
	"os"
	"void/vokernel"
)

var pa string
func InitSocket(){
	pa = RC["socket"]
	println("listening: "+ pa)
	os.RemoveAll(pa)
	Startserver("unix",pa)
}
func Startserver(network string,path string) (*net.Listener,error){
	l,e:=net.Listen(network,path)
	if e!=nil{
		return nil,e
	}else {
		go func() {
			for {
				co, e := l.Accept()
				println("new connection on " + path)
				if e!=nil{
					//log.Print(e)
					break
				}
				_,e=co.Write([]byte("\r\nconnected to void system socket shell\r\n\r\n"))
				if e!=nil{
					log.Print(e)
					break
				}
				go serve(co)
			}
		}()
		return &l, e
	}
}
func serve(co net.Conn){
	rline,wline:=io.Pipe()
	var vw=vokernel.VolatileWriter{Destination: wline}
	go func(){
		io.Copy(&vw,co)
		println("disconnected")
	}()
	sctx:=vokernel.ShellContext{
		WriterSwitch:              &vw, //writer switch
		Reader:                    rline,
		Writer:                    co,
		InternalWriterDestination: wline, //internal receiver
		Delim:                     '\r',
		Privileged:                false,
		Name:                      pa,
	}
	Getsize(sctx)
	for{
		s,e:=sctx.Input(vokernel.Prompt(&sctx))
		if e!=nil{
			println("interrupted")
			sctx.Writer.Close()
			break
		}
		pctx:= PreProcess(s,&sctx)
		Process(pctx)
	}
}