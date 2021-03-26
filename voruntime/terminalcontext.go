package voruntime

import (
	"github.com/chzyer/readline"
	"io"
	"strings"
	"void/vokernel"
)

type Terminal interface {
	Input(prompt string)(string, error)
	Output(content string)
}

type TerminalContext struct {
	RawConnection io.ReadWriteCloser
	StdinReader io.ReadCloser
	StdoutWriter io.WriteCloser
	StdinWriterSwitch *vokernel.VolatileWriter
	internalWriterDestination io.Writer
	Privileged bool
	Delim byte
	ShellName string
	TerminalName string
	runningREPL bool
}
func (t *TerminalContext) RedirectStdinWriter(w io.Writer){
	t.internalWriterDestination=t.StdinWriterSwitch.Destination
	t.StdinWriterSwitch.Destination=w
}
func (t *TerminalContext) RestoreStdinWriter(){
	t.StdinWriterSwitch.Destination=t.internalWriterDestination
}
func (t*TerminalContext) Input(prompt string)(string, error){
	readline.Stdin=t.StdinReader
	readline.Stdout=t.StdoutWriter
	r,_:=readline.New(prompt)
	t.Output("\r")
	return r.Readline()
}
func (t*TerminalContext) InputPassword(prompt string)([]byte, error){
	readline.Stdin=t.StdinReader
	readline.Stdout=t.StdoutWriter
	r,_:=readline.New(prompt)
	t.Output("\r")
	return r.ReadPassword(prompt)
}
func (t*TerminalContext) Output(content string){
	t.StdoutWriter.Write([]byte(strings.ReplaceAll(content,"\n","\r\n")))
}
func (t*TerminalContext) Disconnect(){
	t.RawConnection.Close()
}
func (t *TerminalContext) StartREPL(){
	t.runningREPL=true
	for{
		if !t.runningREPL{
			break
		}
		s,e:=t.Input(Prompt(t))
		if e!=nil{
			println("interrupted")
			t.StdoutWriter.Close()
			break
		}
		pctx:= PreProcess(s,t)
		Process(pctx)
	}
}
func (t *TerminalContext) StopREPL(){
	t.runningREPL=false
	go func(){
		buf:=make([]byte,1)
		for{
			if t.runningREPL{
				break
			}
			t.StdinReader.Read(buf)
			t.StdoutWriter.Write(buf)
		}
	}()
}



func Prompt(tctx *TerminalContext)string{
	if tctx.Privileged{
		return vokernel.Format("<vft green bold>void</vft>:<vft yellow bold>#></vft>")
	}else{
		return vokernel.Format("<vft green bold>void</vft>:<vft blue bold>></vft>")
	}
}
