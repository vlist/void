package voruntime

import (
	"github.com/creack/pty"
	"io"
	exec2 "os/exec"
	"void/vokernel"
)

func BashExec(code string, sctx *vokernel.ShellContext){
	println("exec bash code:"+code)
	args:=[]string{"-c",code}
	Exec(sctx,"bash",args...)
}
func Exec(sctx *vokernel.ShellContext, name string, arg...string){
	rows, cols:=Getsize(*sctx)
	p:=exec2.Command(name,arg...)
	f,_:=pty.StartWithSize(p,&pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
	go io.Copy(sctx.Writer, f)
	sctx.RedirectOutput(f)
	p.Wait()
	sctx.RedirectOutput(sctx.InternalWriterDestination)
}

