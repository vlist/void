package voruntime

import (
	"net"
	"strconv"
	"strings"
	"void/vokernel"
)
type ListenerContext struct{
	Listener *net.Listener
	Flags string
}

var internal map[string]func(*ProcContext)
func InitInternal(){
	internal=map[string]func(*ProcContext){
		"shadow": internal_shadow,
		"info": internal_info,
		"exit": func(pctx *ProcContext) {
			terminal_dispose(pctx.Terminal)
		},
		"shutil": internal_shutil,
		"_stop_repl":func(pctx *ProcContext) {
			pctx.Terminal.StopREPL()
		},
		"su": func(pctx *ProcContext) {
			if !pctx.Terminal.Secured{
				pctx.Terminal.Output("su: Terminal transmission not secured.Reject to authenticate.\n")
				return
			}
			if len(pctx.Args)==0{
				du,e:=Login("guest","guest","")
				if e!=nil{
					pctx.Terminal.Output("su: Login failed: "+e.Error()+"\n")
					return
				}
				pctx.Terminal.User=&du
				return
			}
			user:=strings.Split(pctx.Args[0],":")
			if len(user)==0{
				pctx.Terminal.Output("\nsu: Invalid argument.\n")
				return
			}
			group:=user[0]
			name:=user[1]
			pw,e:=pctx.Terminal.InputPassword("su: Enter password for "+pctx.Args[0]+": ")
			if e!=nil{
				pctx.Terminal.Output("\nsu: Could not login to "+name+".\n")
				return
			}
			lu,e:=Login(name,group,string(pw))
			if e!=nil{
				pctx.Terminal.Output("\nsu: Login failed: "+e.Error()+"\n")
				return
			}
			pctx.Terminal.User=&lu
		},
	}
}
func internal_info(pctx *ProcContext){
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
		pctx.Terminal.Output(vokernel.Format(logo))
	}
	info:= vokernel.GetOSInfo()
	var formattedInfo=""
	formattedInfo+="<vft bold>voidshell</vft> "+info.Version+"\n"
	formattedInfo+="   Runtime/System Arch: "+info.Runtime_SystemArch+"\n"
	pctx.Terminal.Output(vokernel.Format(formattedInfo))
	if printExecContext {
		var formattedExecContext string = ""
		formattedExecContext += "<vft bold>Process Context(pctx):</vft>\n"
		formattedExecContext += "├─ Command Name: " + pctx.CommandName + "\n"
		formattedExecContext += "├─ Arguments: " + "[" + strings.Join(pctx.Args, ",") + "]" + "\n"
		formattedExecContext += "└─ <vft bold>Terminal Context(tctx):</vft>" + "\n"
		formattedExecContext += "   ├─ Shell Interface: " + pctx.Terminal.ShellName + "\n"
		formattedExecContext += "   ├─ Terminal ID: " + pctx.Terminal.TerminalID + "\n"
		formattedExecContext += "   ├─ Transmission Secured: " + strconv.FormatBool(pctx.Terminal.Secured) + "\n"
		formattedExecContext += "   └─ <vft bold>User Context(uctx):</vft>" + "\n"
		formattedExecContext += "      ├─ User Identifier: " + pctx.Terminal.User.Group+":"+pctx.Terminal.User.Name + "\n"
		formattedExecContext += "      └─ Permissions: " + PermissionVisualize(pctx.Terminal.User) + "\n"
		pctx.Terminal.Output(vokernel.Format(formattedExecContext))
	}
}

func terminal_dispose(tctx *TerminalContext){
	disconnectshadow(tctx)
	tctx.Disconnect()
}