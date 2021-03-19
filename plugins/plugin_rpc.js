const net=require("net")
/*
VSRPC
invoke:{
    func: "funcname"
    args: "args"
    id: "id"
}
callback:{
    id: "id"
    result: "result"
}
*/

rpccon=net.connect("/Users/jlywxy/voidshell/rpc1")
rpclist={}
rpccon.on('data',(d)=>{
    rpcobj=JSON.parse(d.toString())
    if(!rpcobj.func){
        rpclist[rpcobj.id](rpcobj.result)
    }else{
        eval(rpcobj.func+"(\""+rpcobj.args+"\")")
    }

})
rpc("rpcready","",(c)=>{
    console.log("rpc callback:" +c)
})
function lrpccall(cmd){
    console.log("process: "+cmd)
}
function rpc(funcname,args,callback){
    id=invoke_uuid()
    rpcobj={
        func: funcname,
        id: id,
        args: args,
        type: "call"
    }
    rpclist[id]=callback
    rpccon.write(JSON.stringify(rpcobj)+"\n")
}
function invoke_uuid() {
    return 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'.replace(/x/g, function (c) {
        var r = Math.random() * 16 | 0,
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

