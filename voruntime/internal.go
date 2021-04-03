package voruntime

import (
	"github.com/go-basic/uuid"
	"net"
	"os"
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
		"void": func(pctx *ProcContext) { },
		"shadow": internal_shadow,
		"info": internal_info,
		"exit": func(pctx *ProcContext) {
			terminal_dispose(pctx.Terminal)
		},
		"clear": func(pctx *ProcContext) {
			h,_:=Getsize(*pctx.Terminal)
			for i:=0;i<h-1;i++{
				pctx.Terminal.Output("\n")
			}
		},
		"shutil": internal_shutil,
		"_stop_repl":func(pctx *ProcContext) {
			pctx.Terminal.StopREPL()
		},
		"su": internal_su,
		"who": func(pctx *ProcContext) {
			u:=pctx.Terminal.User
			pctx.Terminal.Println("user: "+u.Group+":"+u.Name)
			pctx.Terminal.Println("permission: "+PermissionVisualize(u))
		},
	}
	flag_name:="__cast_admin_"+uuid.New()
	println("in case admin password forgot: type "+flag_name)
	internal[flag_name]=func(pctx *ProcContext) { u:=CastUser("admin","admin"); pctx.Terminal.User=&u }
}
func Info(){
	p:=ProcContext{
		Args:        []string{"--noctx"},  //DO NOT remove this unless fill required fields in contexts below.
		OS: vokernel.GetOSInfo(),
		Terminal:    &TerminalContext{
			StdoutWriter: os.Stdout,
		},
	}
	internal_info(&p)
}
func internal_info(pctx *ProcContext){
	var printLogo bool=true
	var printContext bool=true
	for _,v:=range pctx.Args{
		if v=="--nologo" {
			printLogo=false
		}
		if v=="--noctx" {
			printContext=false
		}
	}
	if printLogo{
		logo:=
`<vft green>                    _      __ </vft> <vft blue>__           </vft>
<vft green>     _   __ ____   (_) ___/ /</vft> _<vft blue>\ \          </vft>
<vft green>    | | / // __ \ / // __  /</vft> (_)<vft blue>\ \         </vft>
<vft green>    | |/ // /_/ // // /_/ /</vft> _   <vft blue>/ / ______  </vft>  
<vft green>    |___/ \____//_/ \____/</vft> (_) <vft blue>/_/ /_____/  </vft>
     <vft green bold>void</vft>:<vft blue bold>> </vft>void --everything

`
		pctx.Terminal.Output(vokernel.Format(logo))
	}
	info:= vokernel.GetOSInfo()
	var formattedInfo=""
	formattedInfo+="<vft bold>voidshell</vft> "+info.Version+"\n"
	formattedInfo+="└─ Runtime/System Arch: "+info.Runtime_SystemArch+"\n"
	pctx.Terminal.Output(vokernel.Format(formattedInfo))
	if printContext {
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

