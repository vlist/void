import os,sys
import void
void.info()

class PluginCtx:
    def __init__(self,sctxid,root,self_root):
        self.sctxid=sctxid
        self.root=root+"/root/"
        self.self_root=self_root
    def print(self,content):
        void.print(str(content),self.sctxid)
    def println(self,content):
        void.print(str(content)+"\n",self.sctxid)
    def printf(self,content):
        void.printf(str(content),self.sctxid)
    def input(self,prompt):
        return void.input(prompt,self.sctxid)

def plugin_process(command,sctxid,pluginroot):
    try:
        arg_segs=command.split(" ")
        command_name=arg_segs[0]
        command_args=arg_segs[1:]
        cmdpath=(pluginroot+"/root/"+command_name+".py").replace("//","/")
        cmdpath_index=(pluginroot+"/root/"+command_name+"/__init__.py").replace("//","/")
        ctx=PluginCtx(sctxid,pluginroot,"")
        type="file"
        if os.path.exists(cmdpath) or os.path.exists(cmdpath_index):
            if os.path.exists(cmdpath_index): type="dir"
            cname_segs=command_name.split("/")
            cname_prefix=pluginroot+"/root/"+"/".join(cname_segs[0:len(cname_segs)-1])
            cname=cname_segs[len(cname_segs)-1]
            ctx.self_root=cname_prefix
            if type=="dir":ctx.self_root=cname_prefix+cname+"/"
            try:
                if cname_prefix not in sys.path: sys.path.append(cname_prefix)
                mod=__import__(cname)
            except Exception as e:
                ctx.printf("<vft red bold>[void]</vft>: Error when importing plugin.\n"+str(e)+"\n")
                return
            try:
                mod.init(ctx)
                mod.main(command_args)
            except Exception as e:
                ctx.printf("<vft red bold>[void]</vft>: Error when running plugin.\n"+str(e)+"\n")
                return
            del sys.modules[cname]
        else:
            ctx.printf("<vft red bold>[void]</vft>: command '"+command_name+"' not found.\n")
    except Exception as e:
        print("<vft red bold>[void]</vft>: Error: "+str(e)+"\n")