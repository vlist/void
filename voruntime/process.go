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
	println("process: "+pctx.Type+":"+pctx.CommandName+" "+pctx.Terminal.TerminalID+" "+pctx.Terminal.User.Group+":"+pctx.Terminal.User.Name)
	if pctx.CommandName==""{
		return
	}
	pctx.Terminal.Output("\r")
	switch pctx.Type {
	case "exec":{
		if len(pctx.Args)==0{
			pctx.Terminal.Println("exec: invalid arguments.")
			return
		}
		if pctx.Terminal.User.Permission[1]!=","{
			pctx.Terminal.Println(vokernel.Format("<vft red bold>[void]</vft>: BashExec Permission denied."))
			return
		}
		BashExec(strings.Join(pctx.Args," "),pctx.Terminal)
	}
	case "internal":{
		p,e:=PermissionFilter(pctx.CommandName,pctx.Terminal.User.Permission[0])
		if !p{
			pctx.Terminal.Println(vokernel.Format("<vft red bold>[void]</vft>: Permission denied.\n"+e))
			return
		}
		f := internal[pctx.CommandName]
		if f != nil {
			f(&pctx)
		} else {
			pctx.Terminal.Println("command not found.")
		}}
	case "plugin":{
		p,e:=PermissionFilter(pctx.CommandName,pctx.Terminal.User.Permission[2])
		if !p{
			pctx.Terminal.Println(vokernel.Format("<vft red bold>[void]</vft>: Permission denied.\n"+e))
			return
		}
		//args:=append([]string{"./plugins/plugin_init.js", RC["plugin_root"], pctx.CommandName},pctx.Args...)
		//Exec(pctx.Terminal,"node",args...)
		Plugin_Process(pctx)
		//BashExec("node "+RC["plugin_root"]+"/plugin_init.js "+pctx.CommandName+" "+strings.Join(pctx.Args," "),pctx.Shell)
	}
	}
}
