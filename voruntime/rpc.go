package voruntime

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
)

func InitRPC(){
	var pa="./vsrpc"
	os.RemoveAll(pa)
	l,_:=net.Listen("unix",pa)
	for{
		co,_:=l.Accept()
		go serve(co)
	}
}
type rpcfunc func(string)string
var funclist=map[string]rpcfunc{
	"rpcready": func(string)string{
		println("vs plugin rpc ready")
		return "1"
	},
}
var rpccon net.Conn
var rpclist func(string)string
func serve(co net.Conn){
	rpccon=co
	for{
		r:=bufio.NewReader(co)
		rpcb,e:=r.ReadBytes('\n')
		if e!=nil{
			println("rpc disconnected")
			break
		}
		var rpcobj map[string]string
		json.Unmarshal(rpcb,&rpcobj)
		println("rpc invoke: "+rpcobj["func"]+" "+rpcobj["args"]+" "+rpcobj["id"])
		if rpcobj["func"]!=""{
			f:=funclist[rpcobj["func"]]
			if f!=nil{
				ret:=f(rpcobj["args"])
				println("rpc ret: "+ret)
				var robj= map[string]string{
					"id": rpcobj["id"],
					"result": ret,
				}
				robjb,_:=json.Marshal(&robj)
				co.Write(robjb)
			}
		}
	}
}
func Lrpc(funcname string, args string){
	var rpcobj=map[string]string{
		"func": funcname,
		"args": args,
	}
	robjb,_:=json.Marshal(&rpcobj)
	rpccon.Write(robjb)
}