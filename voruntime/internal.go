package voruntime

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net"
	"os"
	"strings"
	"void/vokernel"
)
type ListenerContext struct{
	Listener *net.Listener
	Flags string
}
var shmap=make(map[string]ListenerContext)
var internal map[string]func(*vokernel.ProcContext)
func InitInternal(){
	internal=map[string]func(*vokernel.ProcContext){

		"info": internal_info,
		"exit": func(pctx *vokernel.ProcContext) {
			pctx.Shell.Writer.Close()
		},
		"sudo": internal_sudo,
		"unsudo": func(pctx *vokernel.ProcContext) {
			pctx.Shell.Privileged=false
		},
		"shutil": func(pctx *vokernel.ProcContext){
			rcsockid:="unix:"+RC["socket"]
			usage:=`usage [--options network:address] [--tls]
options:
	--open: create a new shell socket server
	--kill: close specific socket server
	--list: list all shell socket server
--tls:
	serve over TLS`
			if len(pctx.Args) ==0{
				pctx.Shell.Output(usage)
				return
			}
			switch pctx.Args[0]{
			case "--open":{
				var flag=""
				if len(pctx.Args)<2{
					pctx.Shell.Output("invalid arguments\n")
					pctx.Shell.Output(usage)
					return
				}
				sockid:=pctx.Args[1]
				if sockid==rcsockid{
					pctx.Shell.Output("could not operate on default socket\n")
					return
				}
				na:=strings.Split(sockid,":")
				network:=na[0]
				address:=strings.Join(na[1:],":")
				switch network{
				case "tcp":{}
				case "unix":{
					os.RemoveAll(address)
				}
				default:{
					pctx.Shell.Output("network "+network+" not supported\n")
				}
				}
				var l *net.Listener
				var e error
				//if len(pctx.Args)>=3 && pctx.Args[2]=="--ecdhe-aes"{
				//	println("starting server using ecdhe-aes")
				//	l,e= Startserver_ECDHE_AES(network,address)
				//}else
				if len(pctx.Args)>=3 && pctx.Args[2]=="--tls" {
					flag+="tls "
					println("starting server using TLS")
					l,e= Startserver_TLS(network,address)
				}else{
					l,e= Startserver(network,address)
				}

				if e!=nil{
					pctx.Shell.Output("opening shell on socket "+sockid+" failed\n")
					log.Print(e)
					return
				}
				shmap[sockid]=ListenerContext{
					Listener: l,
					Flags:    flag,
				}
			}
			case "--kill":{
				if len(pctx.Args)<2{
					pctx.Shell.Output("invalid arguments\n")
					pctx.Shell.Output(usage)
					return
				}
				sockid:=pctx.Args[1]
				if sockid==rcsockid{
					pctx.Shell.Output("could not operate on default socket\n")
					return
				}
				if l,ok:=shmap[sockid];ok{
					e:=(*l.Listener).Close()
					if e!=nil{
						pctx.Shell.Output("closing shell on socket "+sockid+" failed\n")
						log.Print(e)
						return
					}
					delete(shmap,sockid)
				}else{
					pctx.Shell.Output("closing shell on socket "+sockid+" failed: listener not found\n")
				}
			}
			case "--list":{
				pctx.Shell.Output("opening socket shell: \n")
				pctx.Shell.Output(rcsockid+"\tdefault\n")
				for k,v := range shmap {
					pctx.Shell.Output(k+"\t"+v.Flags+"\n")
				}
			}
			default:{
				pctx.Shell.Output("invalid arguments\n")
				pctx.Shell.Output(usage)
				return
			}

			}

		},
	}
}
func internal_info(pctx *vokernel.ProcContext){
	var printLogo bool=true
	var printExecContext bool=true
	for _,v:=range pctx.Args{
		if v=="--nologo" {
			printLogo=false
		}
		if v=="--noexeccontext" {
			printExecContext=false
		}
	}
	if printLogo{
		logo:=`
<vft green>                    _      __ </vft> <vft blue>__           </vft>
<vft green>     _   __ ____   (_) ___/ /</vft> _<vft blue>\ \          </vft>
<vft green>    | | / // __ \ / // __  /</vft> (_)<vft blue>\ \         </vft>
<vft green>    | |/ // /_/ // // /_/ /</vft> _   <vft blue>/ / ______  </vft>  
<vft green>    |___/ \____//_/ \____/</vft> (_) <vft blue>/_/ /_____/  </vft>
     <vft green bold>void</vft>:<vft blue bold>></vft>void --everything

`
		pctx.Shell.Output(vokernel.Format(logo))
	}
	info:= vokernel.GetOSInfo()
	var formattedInfo=""
	formattedInfo+="<vft bold>voidshell</vft> "+info.VoVersion+"\n"
	formattedInfo+="    Golang Version: "+info.GoVersion+"\n"
	formattedInfo+="    Current Working Directory: "+info.CurrentWorkingDirectory+"\n"
	formattedInfo+="    System Arch: "+info.SystemArch+"\n"
	pctx.Shell.Output(vokernel.Format(formattedInfo))
	if printExecContext {
		var formattedExecContext string = ""
		formattedExecContext += "<vft bold>Process Context(pctx):</vft>\n"
		formattedExecContext += "    Command Name: " + pctx.CommandName + "\n"
		formattedExecContext += "    Arguments: " + "[" + strings.Join(pctx.Args, ",") + "]" + "\n"
		formattedExecContext += "    Shell Context(sctx): " + "\n"
		formattedExecContext += "        Terminal Name: " + pctx.Shell.Name + "\n"
		formattedExecContext += "        Privileged: " + (func() string {
			if pctx.Shell.Privileged {
				return "true"
			} else {
				return "false"
			}
		})() + "\n"
		pctx.Shell.Output(vokernel.Format(formattedExecContext))
	}
}
func internal_sudo(pctx *vokernel.ProcContext){
	pctx.Shell.Output("sudo: input password\n")
	ipwd,_:=pctx.Shell.InputPassword("")
	h:=sha256.New()
	h.Write(ipwd)
	ipwden:=hex.EncodeToString(h.Sum(nil))
	if RC["password_encrypted"]==ipwden{
		pctx.Shell.Privileged=true
		pctx.Shell.Output("sudo: success\n")
		return
	}else{
		pctx.Shell.Output("sudo: authentication failed\n")
	}
}

