
def init(tctx):
    global ctx
    ctx=tctx

def main(args):
    ctx.print("root: "+ctx.root+", self_root: "+ctx.self_root+"\n")
    ctx.print("testing python plugin_index in go-voidsh😄  (supports UTF-8 if emoji displayed)\n")
    ctx.print("args:")
    ctx.print(str(args)+"\n")
    r=ctx.input("say something? >")
    ctx.print("\r"+r+"\n")
