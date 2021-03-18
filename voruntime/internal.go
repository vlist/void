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

var shmap=map[string]*net.Listener{}
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
			usage:=`usage [--options network:address]
	--open: create a new shell socket server
	--kill: close specific socket server
	--list: list all shell socket server

`
			if len(pctx.Args) ==0{
				pctx.Shell.Output(usage)
				return
			}
			switch pctx.Args[0]{
			case "--open":{
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

				l,e:= Startserver(network,address)
				if e!=nil{
					pctx.Shell.Output("opening shell on socket "+sockid+" failed\n")
					log.Print(e)
					return
				}
				shmap[sockid]=l
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
				l:=shmap[sockid]
				if *l!=nil{
					e:=(*l).Close()
					if e!=nil{
						pctx.Shell.Output("closing shell on socket "+sockid+" failed\n")
						log.Print(e)
						return
					}
					delete(shmap,sockid)
				}
			}
			case "--list":{
				pctx.Shell.Output("opening socket shell: \n")
				pctx.Shell.Output(rcsockid+" (default)\n")
				for k,_ := range(shmap){
					pctx.Shell.Output(k+"\n")
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

