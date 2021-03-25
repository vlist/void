package voruntime

import (
	"github.com/creack/pty"
	"io"
	exec2 "os/exec"
)

func BashExec(code string, tctx *TerminalContext){
	println("exec bash code:"+code)
	args:=[]string{"-c",code}
	Exec(tctx,"bash",args...)
}
func Exec(tctx *TerminalContext, name string, arg...string){
	rows, cols:=Getsize(*tctx)
	p:=exec2.Command(name,arg...)
	f,_:=pty.StartWithSize(p,&pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
	go io.Copy(tctx.StdoutWriter, f)  //process stdout write to terminal stdout
	tctx.RedirectStdinWriter(f)       //terminal stdin write to process stdin
	p.Wait()
	//sctx.RedirectStdinWriter(sctx.InternalWriterDestination)
	tctx.RestoreStdinWriter()
}

