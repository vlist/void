/*
 * voidshell Javascript Plugin Loader
 * Version 1.0
 */
const readline = require('readline');
const vft = require("./vft.js")

_vrs=__dirname.split("/");_vrs.pop()
voidroot=_vrs.join("/")
pluginRoot=process.argv[2]+"/root/"
if(pluginRoot.startsWith(".")){
    //using relative path
    pluginRoot=voidroot+pluginRoot.substr(1)
}

procName=process.argv[3]
procArgv=process.argv.slice(4)
try{
    proc=require(pluginRoot+procName+".js")
}catch(e){
    console.log("voidsh: command '"+procName+"' not found")
    return
}
rl=readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
})
proc.init({
    input: (prompt,callback)=>{
        rl.question(prompt, (answer) => {
            callback(answer)
        });
    },
    print: (content)=>{
        console.log(content)
    },
    printf: (content)=>{
        console.log(vft.format(content))
    },
    format: vft.format,
    exit: ()=>{
        rl.close()
    },
    args: procArgv.unshift(procName),
    root: pluginRoot
})
proc.run()
