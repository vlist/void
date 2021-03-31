package voruntime

import (
	"net"
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
	formattedInfo+="<vft bold>voidshell</vft> "+info.VoVersion+"\n"
	formattedInfo+="    Golang Version: "+info.GoVersion+"\n"
	formattedInfo+="    Current Working Directory: "+info.CurrentWorkingDirectory+"\n"
	formattedInfo+="    System Arch: "+info.SystemArch+"\n"
	pctx.Terminal.Output(vokernel.Format(formattedInfo))
	if printExecContext {
		var formattedExecContext string = ""
		formattedExecContext += "<vft bold>Process Context(pctx):</vft>\n"
		formattedExecContext += "    Command Name: " + pctx.CommandName + "\n"
		formattedExecContext += "    Arguments: " + "[" + strings.Join(pctx.Args, ",") + "]" + "\n"
		formattedExecContext += "    Terminal Context(tctx): " + "\n"
		formattedExecContext += "        Shell Interface: " + pctx.Terminal.ShellName + "\n"
		formattedExecContext += "        Terminal ID: " + pctx.Terminal.TerminalID + "\n"
		pctx.Terminal.Output(vokernel.Format(formattedExecContext))
	}
}

func terminal_dispose(tctx *TerminalContext){
	disconnectshadow(tctx)
	tctx.Disconnect()
}