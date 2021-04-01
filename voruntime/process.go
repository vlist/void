package voruntime

import (
	"strings"
	"void/vokernel"
)

func PreProcess(command string, sctx *TerminalContext) ProcContext {
	segs:=strings.Split(command, " ")
	pctx:= ProcContext{
		CommandName: segs[0],
		Args: segs[1:],
		Terminal: sctx,
	}
	if pctx.CommandName=="exec"{
		pctx.Type="exec"
	}else{
		if internal[pctx.CommandName]!=nil{
			pctx.Type="internal"
		}else{
			pctx.Type="plugin"
		}
	}
	return pctx
}
func Process(pctx ProcContext){
	if strings.TrimSpace(pctx.CommandName)==""{
		return
	}
	println("process: "+pctx.CommandName+" "+pctx.Type+" from "+pctx.Terminal.TerminalID)
	if pctx.CommandName==""{
		return
	}
	pctx.Terminal.Output("\r")
	switch pctx.Type {
	case "exec":{
		if len(pctx.Args)==0{
			pctx.Terminal.Output("exec: invalid arguments\n")
			return
		}
		if pctx.Terminal.User.Permission[1]!=""{
			pctx.Terminal.Output(vokernel.Format("<vft red bold>[void]</vft>: BashExec Permission denied.\n"))
			return
		}
		BashExec(strings.Join(pctx.Args," "),pctx.Terminal)
	}
	case "internal":{
		p,e:=PermissionFilter(pctx.CommandName,pctx.Terminal.User.Permission[0])
		if !p{
			pctx.Terminal.Output(vokernel.Format("<vft red bold>[void]</vft>: Permission denied.\n"+e+"\n"))
			return
		}
		f := internal[pctx.CommandName]
		if f != nil {
			f(&pctx)
		} else {
			pctx.Terminal.Output("command not found\n")
		}}
	case "plugin":{
		p,e:=PermissionFilter(pctx.CommandName,pctx.Terminal.User.Permission[2])
		if !p{
			pctx.Terminal.Output(vokernel.Format("<vft red bold>[void]</vft>: Permission denied.\n"+e+"\n"))
			return
		}
		//args:=append([]string{"./plugins/plugin_init.js", RC["plugin_root"], pctx.CommandName},pctx.Args...)
		//Exec(pctx.Terminal,"node",args...)
		Plugin_Process(pctx)
		//BashExec("node "+RC["plugin_root"]+"/plugin_init.js "+pctx.CommandName+" "+strings.Join(pctx.Args," "),pctx.Shell)
	}
	}
}
