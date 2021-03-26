package voruntime

import (
	"io"
	"void/vokernel"
)
type ShadowState struct{
	srcStdoutWriter io.WriteCloser
	destTerminalName string
	connected bool
}
var shadowstate=make(map[string]ShadowState)
func internal_shadow(pctx *ProcContext) {
	if len(pctx.Args)==0{
		shadow_invalid_argument(pctx.Terminal)
		return
	}
	switch pctx.Args[0] {
	case "-p","--project":{
		if len(pctx.Args)<=1 {
			shadow_invalid_argument(pctx.Terminal)
			return
		}
		disconnectshadow(pctx.Terminal)

		termname:=pctx.Args[1]

		for _,v:=range termmap{
			if v.TerminalName==termname{
				v.Output("shadow connecting to: "+pctx.Terminal.TerminalName+"\n")
				v.Output("--------SHADOW BEGINS--------\n\n")
				v.StdinWriterSwitch.Destination.Write([]byte("_stop_repl\r\n"))
				//go io.Copy(pctx.Terminal.StdinWriterSwitch.Destination,v.StdinReader)

				state:=ShadowState{
					srcStdoutWriter: pctx.Terminal.StdoutWriter,
					destTerminalName: v.TerminalName,
					connected: true,
				}
				shadowstate[pctx.Terminal.TerminalName]=state
				pctx.Terminal.StdoutWriter=vokernel.MultiWriteCloser(v.StdoutWriter,pctx.Terminal.StdoutWriter)
				return
			}
		}
		pctx.Terminal.Output("terminal not found")
	}
	case "-d","--detach":{
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
				tctx.StdoutWriter=state.srcStdoutWriter
				tctx.Output("close existing shadow projector: "+v.TerminalName+"\n")
				v.Output("--------SHADOW ENDS--------\n")
				v.Output("shadow disconnecting from: "+tctx.TerminalName+"\n")
				go v.StartREPL()
				go v.StdinWriterSwitch.Destination.Write([]byte("\r\n"))
				delete(shadowstate,tctx.TerminalName)
				return
			}
		}
	}
}
func shadow_invalid_argument(tctx *TerminalContext){
	tctx.Output("invalid arguments\n")
	usage := `usage [--commands] [terminal name]
commands:
	-p,--project [terminal name]
		project current terminal session to specific terminal
	-d,--detach
		detach shadow terminal
`
	tctx.Output(usage)
}