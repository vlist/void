package voruntime

//func PluginExec(code string, sctx *vokernel.ShellContext){
//	println("plugin code:"+code)
//	p:=exec2.Command("/bin/bash","-c", "node plugin/plugin_init.js "+code)
//	f,_:=pty.Start(p)
//	go func(){
//		for{
//			r:=bufio.NewReader(f)
//			s,e:=r.ReadString('\000')
//			if e!=nil{
//				break;
//			}
//			sctx.Output(vokernel.Format(s))
//		}
//	}()
//	sctx.RedirectOutput(f)
//	p.Wait()
//	sctx.RedirectOutput(sctx.InternalWriterDestination)
//}
