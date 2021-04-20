package voruntime

/*
//BEGIN _VOID_RUNTIME_PLUGIN_CTX_GO_CGO_AUTOFILL_
#cgo pkg-config: python3-embed
//END _VOID_RUNTIME_PLUGIN_CTX_GO_CGO_AUTOFILL_
#cgo CFLAGS: -I. -I${SRCDIR}/..
#include "plugin_def.h"
extern void voidctx_info();
extern void voidctx_print(char*, char*);
typedef PyObject* pyfunc(PyObject*, PyObject*);
*/
import "C"
import (
	"encoding/json"
	"void/vokernel"
)

//export voidctx_info
func voidctx_info(){
	println("* loaded voidshell plugin context. version 1.0")
}
//export voidctx_print_raw
func voidctx_print_raw(content *C.char, tctxid *C.char){
	voidctx_print(C.GoString(content),C.GoString(tctxid))
}
//export voidctx_printf_raw
func voidctx_printf_raw(content *C.char, tctxid *C.char){
	voidctx_printf(C.GoString(content),C.GoString(tctxid))
}
//export voidctx_input_raw
func voidctx_input_raw(content *C.char, tctxid *C.char)*C.char{
	return C.CString(voidctx_input(C.GoString(content),C.GoString(tctxid)))
}
//export voidctx_gettctx_json_raw
func voidctx_gettctx_json_raw(tctxid *C.char)*C.char{
	return C.CString(voidctx_gettctx_json(C.GoString(tctxid)))
}


func voidctx_printf(content string, tctxid string){
	voidctx_print(vokernel.Format(content),tctxid)
}
func voidctx_print(content string,tctxid string){
	if tctx,ok:=termmap[tctxid];ok{
		tctx.Output(content)
	}else{
		println("warning: plugin stdout not piped. ",tctxid)
	}
}
func voidctx_input(prompt string,tctxid string)string{
	if tctx,ok:=termmap[tctxid];ok{
		s,_:=tctx.Input(prompt)
		return s
	}else{
		println("warning: plugin stdin not piped. ",tctxid)
		return ""
	}
}
func voidctx_gettctx_json(tctxid string)string{
	if tctx,ok:=termmap[tctxid];ok{
		s,_:=json.Marshal(tctx)
		return string(s)
	}else{
		println("warning: plugin stdin not piped. ",tctxid)
		return "{}"
	}
}