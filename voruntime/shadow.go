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
			if v.TerminalID==termname{
				v.Println("shadow connecting to: "+pctx.Terminal.TerminalID)
				v.Println("--------SHADOW BEGINS--------\n")
				v.StdinWriterSwitch.Destination.Write([]byte("_stop_repl\r\n"))
				//go io.Copy(pctx.Terminal.StdinWriterSwitch.Destination,v.StdinReader)

				state:=ShadowState{
					srcStdoutWriter: pctx.Terminal.StdoutWriter,
					destTerminalName: v.TerminalID,
					connected: true,
				}
				shadowstate[pctx.Terminal.TerminalID]=state
				pctx.Terminal.StdoutWriter=vokernel.MultiWriteCloser(v.StdoutWriter,pctx.Terminal.StdoutWriter)
				return
			}
		}
		pctx.Terminal.Println("terminal not found.")
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
	if state,ok:=shadowstate[tctx.TerminalID];ok{
		for _,v:=range termmap{
			if v.TerminalID==state.destTerminalName{
				tctx.StdoutWriter=state.srcStdoutWriter
				tctx.Println("close existing shadow projector: "+v.TerminalID)
				v.Println("--------SHADOW ENDS--------")
				v.Println("shadow disconnecting from: "+tctx.TerminalID)
				go v.StartREPL()
				go v.StdinWriterSwitch.Destination.Write([]byte("\r\n"))
				delete(shadowstate,tctx.TerminalID)
				return
			}
		}
	}
}
func shadow_invalid_argument(tctx *TerminalContext){
	tctx.Println("invalid arguments.")
	usage := `usage [--commands] [terminal id]
commands:
	-p,--project [terminal id]
		project current terminal session to specific terminal
	-d,--detach
		detach shadow terminal
`
	tctx.Println(usage)
}