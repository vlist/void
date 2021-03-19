/*
 * Void System Plugin
 * test version 1.0
 */

/*init code,do not modify*/var ctx={};module.exports={ init: (_ctx)=>{ctx=_ctx}, run: main }

function main(){
    ctx.print(ctx.format("<vft green bold>msg from test</vft>"))
    ctx.input("test input>",(t)=>{
        ctx.input("confirm input: "+t+">",(t)=>{
            ctx.print(t)
            ctx.exit()
        })
    })
}
