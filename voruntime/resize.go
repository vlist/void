package voruntime

import (
	"bufio"
	"io"
	"strconv"
)

func Getsize(tctx TerminalContext)(int,int){
	r,w:=io.Pipe()
	//checkout https://unix.stackexchange.com/questions/16578/resizable-serial-console-window
	tctx.Output("\0337\033[r\033[999;999H\033[6n\0338")
	tctx.RedirectStdinWriter(w)
	f:=bufio.NewReader(r)
	f.ReadBytes('[')
	colsB,_:=f.ReadBytes(';')
	rowsB,_:=f.ReadBytes('R')
	//sctx.RedirectStdinWriter(sctx.InternalWriterDestination)
	tctx.RestoreStdinWriter()
	rows,_:=strconv.Atoi(string(colsB[0:len(colsB)-1]))
	cols,_:=strconv.Atoi(string(rowsB[0:len(rowsB)-1]))
	return rows,cols
}

