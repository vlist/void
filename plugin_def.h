#ifndef _VOID_PLUGIN_DEF_H_
#define _VOID_PLUGIN_DEF_H_

#define PY_SSIZE_T_CLEAN
#include <Python.h>

void voidctx_info();
void voidctx_print_raw(char*, char*);
void voidctx_printf_raw(char*, char*);
char* voidctx_input_raw(char*, char*);
char* voidctx_gettctx_json_raw(char*);

#endif
