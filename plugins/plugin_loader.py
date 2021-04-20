import os,sys,json
import void
void.info()

class PluginCtx:
    def __init__(self,tctxid,root,self_root):
        self.tctxid=tctxid
        try:
            self.tctx=json.loads(void.get_tctx_json(tctxid))
        except:
            pass
        self.root=root+"/root/"
        self.self_root=self_root
    def print(self,content):
        void.print(str(content),self.tctxid)
    def println(self,content):
        void.print(str(content)+"\n",self.tctxid)
    def printf(self,content):
        void.printf(str(content),self.tctxid)
    def input(self,prompt):
        return void.input(prompt,self.tctxid)

def plugin_process(command,tctxid,pluginroot):
    try:
        arg_segs=command.split(" ")
        command_name=arg_segs[0]
        command_args=arg_segs[1:]
        cmdpath=(pluginroot+"/root/"+command_name+".py").replace("//","/")
        cmdpath_index=(pluginroot+"/root/"+command_name+"/__init__.py").replace("//","/")
        ctx=PluginCtx(tctxid,pluginroot,"")
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
                if cname in sys.modules: del sys.modules[cname]
                mod=__import__(cname)
            except Exception as e:
                ctx.printf("<vft red bold>[void]</vft>: Error when importing plugin.\n"+str(e)+"\n")
                return
            try:
                mod.init(ctx)
                r=mod.main(command_args)
                if r==None: r=0
                color=""
                if r!=0: color="red"
                ctx.printf("\n<vft " +color+ " bold>[void]</vft>: '"+cname+"' exit "+str(r)+".\n")
            except Exception as e:
                ctx.printf("<vft red bold>[void]</vft>: Error when running plugin.\n"+str(e)+"\n")
                return
        else:
            ctx.printf("<vft red bold>[void]</vft>: command '"+command_name+"' not found.\n")
    except Exception as e:
        print("<vft red bold>[void]</vft>: Error: "+str(e)+"\n")