package vokernel

import (
	"github.com/chzyer/readline"
	"io"
	"strings"
)

type ShellContext struct {
	Reader io.ReadCloser
	Writer io.WriteCloser
	WriterSwitch *VolatileWriter
	InternalWriterDestination io.Writer
	Width uint16
	Height uint16
	Privileged bool
	Delim byte
	Name string
}
func (t *ShellContext) RedirectOutput(w io.Writer){
	t.WriterSwitch.Destination=w
}
func (t* ShellContext) Input(prompt string)(string, error){
	readline.Stdin=t.Reader
	readline.Stdout=t.Writer
	r,_:=readline.New(prompt)
	t.Output("\r")
	return r.Readline()
}
func (t* ShellContext) InputPassword(prompt string)([]byte, error){
	readline.Stdin=t.Reader
	readline.Stdout=t.Writer
	r,_:=readline.New(prompt)
	t.Output("\r")
	return r.ReadPassword(prompt)
}
func (t* ShellContext) Output(content string){
	t.Writer.Write([]byte(strings.ReplaceAll(content,"\n","\r\n")))
}
func Prompt(sctx *ShellContext)string{
	if sctx.Privileged{
		return Format("<vft green bold>void</vft>:<vft yellow bold>#></vft>")
	}else{
		return Format("<vft green bold>void</vft>:<vft blue bold>></vft>")
	}
}