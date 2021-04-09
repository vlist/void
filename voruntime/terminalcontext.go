package voruntime

import (
	"encoding/hex"
	"github.com/chzyer/readline"
	"io"
	"strings"
	"void/vokernel"
)

type Terminal interface {
	Input(prompt string)(string, error)
	Output(content string)
}

var termmap=make(map[string]*TerminalContext)

type TerminalContext struct {
	RawConnection io.Closer
	StdinReader io.ReadCloser
	StdoutWriter io.WriteCloser
	StdinWriterSwitch *vokernel.VolatileWriter
	internalWriterDestination io.Writer
	Secured bool
	Delim byte
	ShellName string
	TerminalID string
	runningREPL bool
	User *UserContext
	Environment map[string]interface{}
}
func (t *TerminalContext) RedirectStdinWriter(w io.Writer){
	t.internalWriterDestination=t.StdinWriterSwitch.Destination
	t.StdinWriterSwitch.Destination=w
}
func (t *TerminalContext) RestoreStdinWriter(){
	t.StdinWriterSwitch.Destination=t.internalWriterDestination
}
type KeyListener struct{
	Terminal *TerminalContext
}
func (l * KeyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool){
	print(hex.EncodeToString([]byte(string(key))))
	return line,pos,true
}
func (t*TerminalContext) Input(prompt string)(string, error){
	readline.Stdin=t.StdinReader
	readline.Stdout=t.StdoutWriter
	//var l readline.Listener= &KeyListener{t}
	r,_:=readline.NewEx(&readline.Config{
		Prompt:                 prompt,
		HistoryFile:            ".voidsh_history/"+t.TerminalID,
		HistoryLimit:           100,
		//Listener:               l,
	})
	t.Output("\r")
	line,e:=r.Readline()
	r.Close()
	return line,e
}
func (t*TerminalContext) InputPassword(prompt string)([]byte, error){
	readline.Stdin=t.StdinReader
	readline.Stdout=t.StdoutWriter
	r,_:=readline.New(prompt)
	l,e:=r.ReadPassword(prompt)
	t.Output("\r")
	r.Close()
	return l,e
}
func (t*TerminalContext) Output(content string){
	t.StdoutWriter.Write([]byte(strings.ReplaceAll(content,"\n","\r\n")))
}
func (t*TerminalContext) Println(content string){
	t.Output(content+"\n")
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

func clientHello(tctx *TerminalContext){
	pctx:=ProcContext{
		CommandName: "",
		Args:        []string{},
		Type:        "",
		Terminal:    tctx,
		OS:          vokernel.OSInfo{},
	}
	internal_info(&pctx)
}

func Prompt(tctx *TerminalContext)string{
	if tctx.User.Group=="guest"{
		return vokernel.Format("<vft green bold>void</vft>:<vft blue bold>></vft> ")
	}else{
		return vokernel.Format("<vft green bold>void</vft>:<vft yellow bold>"+tctx.User.Name+"#</vft> ")
	}
}
func terminal_dispose(tctx *TerminalContext){
	//os.RemoveAll(".voidsh_history/"+tctx.TerminalID)
	disconnectshadow(tctx)
	tctx.Disconnect()
}