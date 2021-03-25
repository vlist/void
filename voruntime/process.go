package voruntime

import (
	"strings"
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
	println("process: "+pctx.CommandName+" "+pctx.Type+" from "+pctx.Terminal.TerminalName)
	if pctx.CommandName==""{
		return
	}
	pctx.Terminal.Output("\r")
	switch pctx.Type {
	case "exec":{
		BashExec(strings.Join(pctx.Args," "),pctx.Terminal)
	}
	case "internal":{
			f := internal[pctx.CommandName]
			if f != nil {
				f(&pctx)
			} else {
				pctx.Terminal.Output("command not found\n")
			}}
	case "plugin":{
		args:=append([]string{"./plugins/plugin_init.js", RC["plugin_root"], pctx.CommandName},pctx.Args...)
		Exec(pctx.Terminal,"node",args...)
		//BashExec("node "+RC["plugin_root"]+"/plugin_init.js "+pctx.CommandName+" "+strings.Join(pctx.Args," "),pctx.Shell)
	}
	}
}
