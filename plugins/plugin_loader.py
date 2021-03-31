import os,sys
import void
from importlib import reload
void.info()

class PluginCtx:
    def __init__(self,sctxid):
        self.sctxid=sctxid
    def print(self,content):
        void.print(str(content),self.sctxid)
    def printf(self,content):
        void.printf(str(content),self.sctxid)
    def input(self,prompt):
        return void.input(prompt,self.sctxid)

def plugin_process(command,sctxid):
    arg_segs=command.split(" ")
    command_name=arg_segs[0]
    command_args=arg_segs[1:]
    cmdpath=("./plugins/root/"+command_name+".py").replace("//","/")
    ctx=PluginCtx(sctxid)
    if os.path.exists(cmdpath):
        cname_segs=command_name.split("/")
        cname_prefix="./plugins/root/"+"/".join(cname_segs[0:len(cname_segs)-1])
        cname=cname_segs[len(cname_segs)-1]
        try:
            if cname_prefix not in sys.path: sys.path.append(cname_prefix)
            mod=__import__(cname)
        except Exception as e:
            ctx.printf("<vft red bold>[void]</vft>: Error when importing plugin.\n"+str(e)+"\n")
        try:
            mod.init(ctx)
            mod.main(command_args)
        except Exception as e:
            ctx.printf("<vft red bold>[void]</vft>: Error when trying to init and run plugin.\n"+str(e)+"\n")
        del sys.modules[cname]
    else:
        ctx.printf("<vft red bold>[void]</vft>: command '"+command_name+"' not found.\n")