#!/bin/bash
# What will this script do:
# 1. check whether python-dev(el) has installed using "python3-config"
# 2. find pkg-config file (python3.pc) under python3-config prefix
#    (may take a lot of time)
# 3. get and set PKG_CONFIG_PATH from pkg-config file path
# 4. determine to use "python3" or "python3-embed" of pkg-config to cgo using "pkg-config"
echo
echo CGO-CPytion Configure Tool for void \(1.0 20R14a\)
echo
echo Re-run this tool whenever python has updated or build target has changed.

PY_PREFIX=`python3-config --prefix`
if [ $? != 0 ]; then
    echo \* checking python3-config prefix... no
    echo Could not determine where is python3 prefix. Have python3-dev installed?
    exit 1
else
    echo \* checking python3-config prefix... ok
fi
PY_PC_FILE=`find $PY_PREFIX -name python3.pc`
if [ "$PY_PC_FILE" == "" ]; then
    echo \* finding python3 pkg-config file... no
    echo Could not find python3.pc file. Have python3-dev installed?
    exit 1
else
    echo \* finding python3 pkg-config file... ok
fi
PCP_RAW=${PY_PC_FILE%/*}
echo \* setting PKG_CONFIG_PATH: $PCP_RAW
export PKG_CONFIG_PATH=$PCP_RAW:$PKG_CONFIG_PATH
pkg-config --list-all|grep python3-embed >> /dev/null
if [ $? != 0 ]; then
    pkg-config --list-all|grep python3>> /dev/null
    echo \* checking pkg-config python3... ok
    export VO_BUILD_CGO="#cgo pkg-config: python3"
else
    echo \* checking pkg-config python3-embed... ok
    export VO_BUILD_CGO="#cgo pkg-config: python3-embed"
fi
echo ------------------------------
echo "(if needed) Set ⬇️ in .bash_profile or Goland Run/Debug Configuration -> Environment:"
echo PKG_CONFIG_PATH=$PCP_RAW
echo
echo "Put ⬇️ in plugin.go file:"
echo $VO_BUILD_CGO
echo ------------------------------











# pkg-config --list-all|grep python >> /dev/null
# if [ $? != 0 ]; then
#     echo \* checking python3 pkg-config environment... no
#     echo
#     echo \# Put in .bash_profile or equivalent file or set environment in Goland Run/Debug Configurations:
#     echo ----------------------------------------------------
#     echo export PKG_CONFIG_PATH=${PY_PC_FILE%/*}:\$PKG_CONFIG_PATH
#     echo ----------------------------------------------------
#     echo
#     export PKG_CONFIG_PATH=${PY_PC_FILE%/*}:\$PKG_CONFIG_PATH
# else
#     echo \* checking python3 pkg-config environment... ok
# fi
# export PKG_CONFIG_PATH_PY3=${PY_PC_FILE%/*}
# python3 - "$@" <<END
# #!/usr/bin/python3
# import os
# f=open(os.environ["PKG_CONFIG_PATH_PY3"]+"/python3.pc")
# if f.read().find("Libs:\n")!=-1:
#     exit(1)
# exit(0)
# END

# if [ $? != 0 ]; then
#     echo \* checking python3 pkg-config LDFLAGS... no
#     PKG_CONFIG_PATH_PY3_NO_LDFLAGS=1
#     # LDFLAGS not found in python3.pc, generating from `python3-config ...`
# else
#     echo \* checking python3 pkg-config LDFLAGS... ok
# fi
# echo
# echo // Put in go source file: 
# echo ----------------------------------------------------
# echo \#cgo pkg-config: python3
# if [ "$PKG_CONFIG_PATH_PY3_NO_LDFLAGS" == "1" ]; then
#     PY_LD_FLAG=`python3-config --ldflags --embed`
#     if [ $? != 0 ]; then
#         PY_LD_FLAG=`python3-config --ldflags`
#     fi
#     echo \#cgo LDFLAGS: $PY_LD_FLAG
#     export PY_LD_FLAG
# fi
# echo ----------------------------------------------------