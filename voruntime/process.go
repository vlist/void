package voruntime

import (
	"strconv"
	"strings"
	"void/vokernel"
)

func PreProcess(command string, sctx *vokernel.ShellContext) vokernel.ProcContext{
	segs:=strings.Split(command, " ")
	pctx:=vokernel.ProcContext{
		CommandName: segs[0],
		Args: segs[1:],
		Shell: sctx,
	}
	if pctx.CommandName=="setsize"{
		w,_:=strconv.Atoi(segs[2])
		h,_:=strconv.Atoi(segs[1])
		pctx.Shell.Width= uint16(int(w))
		pctx.Shell.Width= uint16(int(h))
	}else if pctx.CommandName=="exec"{
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
func Process(pctx vokernel.ProcContext){
	println("process: "+pctx.CommandName)
	if pctx.CommandName==""{
		return
	}
	pctx.Shell.Output("\r")
	switch pctx.Type {
	case "exec":{
		Exec(strings.Join(pctx.Args," "),pctx.Shell)
	}
	case "internal":{
			f := internal[pctx.CommandName]
			if f != nil {
				f(&pctx)
			} else {
				pctx.Shell.Output("command not found\n")
			}}
	case "plugin":{
		Exec("node plugin/plugin_init.js "+pctx.CommandName+" "+strings.Join(pctx.Args," "),pctx.Shell)
	}
	}
}
