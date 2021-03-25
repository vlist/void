package voruntime

import (
	"io"
	"void/vokernel"
)
type ShadowState struct{
	srcStdoutWriter io.WriteCloser
	destTerminalName string
}
var shadowstate=make(map[string]ShadowState)
func internal_shadow(pctx *ProcContext) {
	switch pctx.Args[0] {
	case "--project":{
		if !(len(pctx.Args)>=2) {
			shadow_invalid_argument(pctx.Terminal)
		}
		disconnectshadow(pctx.Terminal)

		termname:=pctx.Args[1]

		for _,v:=range termmap{
			if v.TerminalName==termname{
				v.Output("shadow connecting to: "+pctx.Terminal.TerminalName+"\n")
				v.Output("--------SHADOW BEGINS--------\n\n")
				v.StdinWriterSwitch.Destination.Write([]byte("_stop_repl\r\n\r\n"))

				state:=ShadowState{
					srcStdoutWriter: pctx.Terminal.StdoutWriter,
					destTerminalName: v.TerminalName,
				}
				shadowstate[pctx.Terminal.TerminalName]=state
				//pctx.Terminal.StdoutWriter=vokernel.BiWriteCloser(pctx.Terminal.StdoutWriter,v.StdoutWriter)
				pctx.Terminal.StdoutWriter=vokernel.MultiWriteCloser(v.StdoutWriter,pctx.Terminal.StdoutWriter)
				return
			}
		}
		pctx.Terminal.Output("terminal not found")
	}
	case "--detach":{
		disconnectshadow(pctx.Terminal)
	}
	default:{
		shadow_invalid_argument(pctx.Terminal)
	}
	}

}
func disconnectshadow(tctx *TerminalContext){
	if state,ok:=shadowstate[tctx.TerminalName];ok{
		for _,v:=range termmap{
			if v.TerminalName==state.destTerminalName{
				tctx.Output("close existing shadow projector: "+v.TerminalName+"\n")
				v.Output("\n\n--------SHADOW ENDS--------\n")
				v.Output("shadow disconnecting from: "+tctx.TerminalName+"\n")
				tctx.StdoutWriter=state.srcStdoutWriter
				v.StartREPL()
				v.StdinWriterSwitch.Destination.Write([]byte("\r\n\r\n"))
				return
			}
		}
	}
	delete(shadowstate,tctx.TerminalName)
}
func shadow_invalid_argument(tctx *TerminalContext){
	tctx.Output("invalid argument")
	usage := `usage [commands] [terminal name]
commands:
	--attach
		attach current terminal session to specific terminal
	--detach
		detach current terminal from specific terminal
	--switch
		reattach current terminal session to specific terminal
`
	tctx.Output(usage)
}