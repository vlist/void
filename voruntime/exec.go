package voruntime

import (
	"github.com/creack/pty"
	"io"
	exec2 "os/exec"
	"void/vokernel"
)

func Exec(code string, sctx *vokernel.ShellContext){
	println("exec code:"+code)
	bfix:="resize>>/dev/null"
	efix:=bfix
	p:=exec2.Command("/bin/bash","-c", bfix+";"+code+";"+efix)
	f,_:=pty.Start(p)
	go io.Copy(sctx.Writer, f)
	sctx.RedirectOutput(f)
	p.Wait()
	sctx.RedirectOutput(sctx.InternalWriterDestination)
}
