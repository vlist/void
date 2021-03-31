package voruntime

import (
	"crypto/tls"
	"github.com/go-basic/uuid"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"void/vokernel"
)



func InitSocket(){
	pa := RC["socket"]
	println("listening: "+ pa)
	os.RemoveAll(pa)
	dir:=filepath.Dir(pa)
	err:=os.MkdirAll(dir,0770)
	if err!=nil{
		println("checking directory failed: "+dir)
		return
	}
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
				_,e=co.Write([]byte("\r\nconnected to voidshell\r\n\r\n"))
				if e!=nil{
					log.Print(e)
					break
				}
				go serve(co,network+":"+path)
			}
		}()
		return &l, e
	}
}
func Startserver_TLS(network string,path string) (*net.Listener,error){
	cert,e:=tls.LoadX509KeyPair(RC["tls_config_pem"],RC["tls_config_key"])//"cert/server.pem","cert/server.key")
	cfg:=&tls.Config{Certificates: []tls.Certificate{cert}}
	l,e:=tls.Listen(network,path,cfg)
	if e!=nil{
		return nil,e
	}else {
		go func() {
			for {
				co, e := l.Accept()
				println("new connection on " + path)
				if e!=nil{
					println("err when accept")
					log.Print(e)
					break
				}
				_,e=co.Write([]byte("\r\nconnected to voidshell\r\n\r\n"))
				if e!=nil{
					println("err when write")
					log.Print(e)
					//break
				}
				go serve(co,"tls:"+path)
			}
			println("tls accept stopped")
		}()
		return &l, e
	}
}
func serve(co net.Conn,servername string){
	stdinReader, socketStdinWriter:=io.Pipe()
	termid:=uuid.New()
	var stdinWriterVolatile=vokernel.VolatileWriter{Destination: socketStdinWriter}
	go func(){
		io.Copy(&stdinWriterVolatile,co)  //socket write to stdin writer, shell read from stdin reader
		println("disconnected")
		delete(termmap,termid)
	}()

	tctx:= TerminalContext{
		RawConnection:					co,
		StdinWriterSwitch:              &stdinWriterVolatile,
		StdinReader:                    stdinReader,
		StdoutWriter:                   co,
		Delim:                     		'\r',
		Privileged:                		false,
		ShellName:                      servername,
		TerminalID: 					termid,
	}
	termmap[termid]=&tctx
	clientHello(&tctx)
	go tctx.StartREPL()
}
//func Startserver_ECDHE_AES(network string,path string) (*net.Listener,error){
//	l,e:=net.Listen(network,path)
//	if e!=nil{
//		return nil,e
//	}else {
//		go func() {
//			for {
//				co, e := l.Accept()
//				println("new connection (aes cipher) on " + path)
//				if e!=nil{
//					break
//				}
//				go serve_aes(co,network+":"+path)
//			}
//		}()
//		return &l, e
//	}
//}
//func serve_aes(raw_co net.Conn,servername string){
//	p256 := ecdh.Generic(secp256k1.S256())
//	privkey,pubkey,_:=p256.GenerateKey(rand.Reader)
//	ppubkey:=pubkey.(ecdh.Point)
//	pub_b:=elliptic.Marshal(secp256k1.S256(),ppubkey.X,ppubkey.Y)
//	println("this pubkey "+hex.EncodeToString(pub_b))
//
//	cli_pub_b:=make([]byte,65)
//	raw_co.Read(cli_pub_b)
//	println("get client pubkey "+hex.EncodeToString(cli_pub_b))
//	x,y:=elliptic.Unmarshal(secp256k1.S256(),cli_pub_b)
//	if x==nil||y==nil{
//		println("pubkey invalid")
//		raw_co.Close()
//		return
//	}
//	cli_pub:=ecdh.Point{x,y}
//	raw_co.Write(pub_b)
//	key:=p256.ComputeSecret(privkey,cli_pub)
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		panic(err)
//	}
//	var iv [aes.BlockSize]byte
//	stream := cipher.NewCTR(block, iv[:])
//	sec_co_in:=cipher.StreamReader{
//		S: stream,
//		R: raw_co,
//	}
//	sec_co_out:=cipher.StreamWriter{
//		S: stream,
//		W: raw_co,
//	}
//	sec_co_out.Write([]byte("\r\nconnected to void system socket shell\r\n\r\n"))
//	var co = vokernel.VolatileReader{Source: sec_co_in}
//	rline,wline:=io.Pipe()
//	var vw=vokernel.VolatileWriter{Destination: wline}
//	go func(){
//		io.Copy(&vw, &co)
//		println("disconnected")
//	}()
//	sctx:=vokernel.ShellContext{
//		WriterSwitch:              &vw, //writer switch
//		Reader:                    rline,
//		Writer:                    sec_co_out,
//		InternalWriterDestination: wline, //internal receiver
//		Delim:                     '\r',
//		Privileged:                false,
//		Name:                      servername,
//	}
//	Getsize(sctx)
//	for{
//		s,e:=sctx.Input(vokernel.Prompt(&sctx))
//		if e!=nil{
//			println("interrupted")
//			sctx.Writer.Close()
//			break
//		}
//		pctx:= PreProcess(s,&sctx)
//		Process(pctx)
//	}
//}
//
