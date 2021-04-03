package voruntime
// IMPORTANT:
/*
1. DO NOT move(project relative path) or rename this file
unless modify filename in build.sh line 13 and 22
and ⬇️ #cgo CFLAGS -I. -I${SRCDIR}/..
${SRCDIR} is representing to path of this .go source file.

2. DO NOT modify tag ⬇️ ️_VOID_RUNTIME_PLUGIN_GO_CGO_AUTOFILL_
if you want to use build.sh
*/

/*
//BEGIN _VOID_RUNTIME_PLUGIN_GO_CGO_AUTOFILL_
#cgo pkg-config: python3-embed
//END _VOID_RUNTIME_PLUGIN_GO_CGO_AUTOFILL_
#cgo CFLAGS: -I. -I${SRCDIR}/..// -I"/Library/Java/JavaVirtualMachines/jdk1.8.0_231.jdk/Contents/Home/include" -I"/Library/Java/JavaVirtualMachines/jdk1.8.0_231.jdk/Contents/Home/include/darwin"
#include "plugin_def.h"

PyObject* func_voidctx_info(PyObject *self, PyObject *args) {
	voidctx_info();
	return PyLong_FromLong(0);
}

//check out https://stackoverflow.com/questions/62413050/c-char-array-from-python-string
char* str(PyObject *o){
	return PyBytes_AsString(PyUnicode_AsEncodedString(o,"UTF-8", "strict"));
}
PyObject* func_voidctx_print_raw(PyObject *self, PyObject *args) {
	voidctx_print_raw(str(PyTuple_GetItem(args,0)),str(PyTuple_GetItem(args,1)));
	return PyLong_FromLong(0);
}
PyObject* func_voidctx_printf_raw(PyObject *self, PyObject *args) {
	voidctx_printf_raw(str(PyTuple_GetItem(args,0)),str(PyTuple_GetItem(args,1)));
	return PyLong_FromLong(0);
}
PyObject* func_voidctx_input_raw(PyObject *self, PyObject *args) {
	return PyUnicode_FromString(voidctx_input_raw(str(PyTuple_GetItem(args,0)),str(PyTuple_GetItem(args,1))));
}
PyMethodDef voidctxMethods[] = {
	{"info", func_voidctx_info, METH_VARARGS, NULL},
	{"print", func_voidctx_print_raw, METH_VARARGS, NULL},
	{"printf", func_voidctx_printf_raw, METH_VARARGS, NULL},
	{"input", func_voidctx_input_raw, METH_VARARGS, NULL},
	{NULL, NULL, 0, NULL}
};
PyModuleDef voidctxModule = {
	PyModuleDef_HEAD_INIT, "void", NULL, -1, voidctxMethods,
	NULL, NULL, NULL, NULL
};
PyObject* PyInit_voidctx(void) { return PyModule_Create(&voidctxModule); }
void voidctxInit(){
	PyImport_AppendInittab("void", &PyInit_voidctx);
}

*/
import "C"
import (
	"strings"
	"unsafe"
	"void/vokernel"
)

var voidplugin_process *C.PyObject
var plugin_func_inited bool

func InitPlugin() {
	plugin_func_inited=false
	C.voidctxInit()
	C.Py_Initialize()
	initcode:=C.CString("import sys;sys.path.append(\"./plugins\")")
	C.PyRun_SimpleString(initcode)
	C.free(unsafe.Pointer(initcode))

	path:=C.CString("plugin_loader")
	loader_mod:=C.PyImport_ImportModule(path)
	C.free(unsafe.Pointer(path))
	if loader_mod==nil{
		println("could not import plugin loader")
		return
	}
	process_fname:=C.CString("plugin_process")
	voidplugin_process=C.PyObject_GetAttrString(loader_mod,process_fname)
	C.free(unsafe.Pointer(process_fname))
	if voidplugin_process==nil{
		println("could not access plugin process function")
		return
	}
	plugin_func_inited=true
}

func Plugin_Process(pctx ProcContext){
	if !plugin_func_inited{
		pctx.Terminal.Println("voidshell plugin loader mod could not be initialized. Check syntax error in plugin_loader.py")
		return
	}
	args:=C.PyTuple_New(3)
	arg0_raw:=C.CString(pctx.CommandName+" "+strings.Join(pctx.Args," "))
	arg1_raw:=C.CString(pctx.Terminal.TerminalID)
	arg2_raw:=C.CString(vokernel.RC["plugin_root"])
	arg0:=C.PyUnicode_FromString(arg0_raw)
	arg1:=C.PyUnicode_FromString(arg1_raw)
	arg2:=C.PyUnicode_FromString(arg2_raw)
	C.free(unsafe.Pointer(arg0_raw))
	C.free(unsafe.Pointer(arg1_raw))
	C.free(unsafe.Pointer(arg2_raw))
	C.PyTuple_SetItem(args,0,arg0);C.PyTuple_SetItem(args,1,arg1);C.PyTuple_SetItem(args,2,arg2)
	if ret:=C.PyObject_CallObject(voidplugin_process,args);ret==nil{
		pctx.Terminal.Println(vokernel.Format("<vft red bold>[void]</vft>: Could not execute plugin."))
		return
	}
}