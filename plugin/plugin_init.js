/*
 * Void System Javascript Plugin Loader
 * Version 1.0
 */
const readline = require('readline');
const vft = require("./vft.js")

procName=process.argv[2]
procArgv=process.argv.slice(3)
try{
    proc=require(__dirname+"/root/"+procName+".js")
}catch(e){
    console.log("voidsh: command '"+procName+"' not found")
    return
}
rl=readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
})
/*
 * Plugin Context
 * format: Void Format Text
 * input: readline.question
 * output: console.log(...)
 * exit: call before plugin terminated(otherwise send SIGINT by user to terminate)
 * args: [pluginname, plugin args...]
 */
proc.init({
    input: (prompt,callback)=>{
        rl.question(prompt, (answer) => {
            callback(answer)
        });
    },
    print: (content)=>{
        console.log(content)
    },
    format: vft.format,
    exit: ()=>{
        rl.close()
    },
    args: procArgv.unshift(procName)
})
proc.run()
