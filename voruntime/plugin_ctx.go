package voruntime

/*
#cgo CFLAGS: -I. -I${SRCDIR}/..
#include "plugin_def.h"
extern void voidctx_info();
extern void voidctx_print(char*, char*);
*/
import "C"
import "void/vokernel"

//export voidctx_info
func voidctx_info(){
	println("void plugin context v0.1")
}
//export voidctx_print_raw
func voidctx_print_raw(content *C.char, sctxid *C.char){
	voidctx_print(C.GoString(content),C.GoString(sctxid))
}
//export voidctx_printf_raw
func voidctx_printf_raw(content *C.char, sctxid *C.char){
	voidctx_printf(C.GoString(content),C.GoString(sctxid))
}
//export voidctx_input_raw
func voidctx_input_raw(content *C.char, sctxid *C.char)*C.char{
	return C.CString(voidctx_input(C.GoString(content),C.GoString(sctxid)))
}

func voidctx_printf(content string, sctxid string){
	voidctx_print(vokernel.Format(content),sctxid)
}
func voidctx_print(content string,sctxid string){
	if tctx,ok:=termmap[sctxid];ok{
		tctx.Output(content)
	}else{
		println("warning: plugin stdout not piped. ",sctxid)
	}
}
func voidctx_input(prompt string,sctxid string)string{
	if tctx,ok:=termmap[sctxid];ok{
		s,_:=tctx.Input(prompt)
		return s
	}else{
		println("warning: plugin stdin not piped. ",sctxid)
		return ""
	}
}