package voruntime

import (
	"github.com/creack/pty"
	"io"
	exec2 "os/exec"
	"void/vokernel"
)

func Exec(code string, sctx *vokernel.ShellContext){
	println("exec code:"+code)
	//bfix:="IFS=\" \";read -a size <<< `stty size`;stty rows ${size[0]} cols ${size[1]};echo rows=${size[0]},cols=${size[1]}"
	rows, cols:=Getsize(*sctx)
	//bfix:="stty rows "+strconv.Itoa(rows)+" cols "+strconv.Itoa(cols)
	//bfix:="echo"
	//efix:=bfix
	p:=exec2.Command("/bin/bash","-c", code)
	f,_:=pty.StartWithSize(p,&pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
	go io.Copy(sctx.Writer, f)
	sctx.RedirectOutput(f)
	p.Wait()
	sctx.RedirectOutput(sctx.InternalWriterDestination)
}
