import json
def init(tctx):
    global ctx
    ctx=tctx

def main(args):
    ctx.println("voidshell Plugin Extension 1.0")
    ctx.println("Terminal Context: ")
    ctx.println(ctx.__dict__)
    return 0