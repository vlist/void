package voruntime

import (
	"strings"
	"void/vokernel"
)

func internal_su(pctx *ProcContext) {
	if !pctx.Terminal.Secured {
		if vokernel.RC["allow_su_via_insecure_transmission"] =="false"{
			pctx.Terminal.Println("su: Error: Terminal transmission not secured. Reject to authenticate.")
			return
		}else {
			pctx.Terminal.Println("su: Warning: Terminal transmission not secured. Typing administrator passwords is not recommended.")
		}
	}
	if len(pctx.Args)==0{
		du,e:=Login("guest","guest","")
		if e!=nil{
			pctx.Terminal.Println("su: Login failed: "+e.Error())
			return
		}
		pctx.Terminal.User=&du
		return
	}
	user:=strings.Split(pctx.Args[0],":")
	if len(user)<2{
		pctx.Terminal.Println("\nsu: Invalid argument.")
		return
	}
	group:=user[0]
	name:=user[1]
	pw,e:=pctx.Terminal.InputPassword("su: Enter password for "+pctx.Args[0]+": ")
	if e!=nil{
		pctx.Terminal.Println("\nsu: Could not login to "+name+".")
		return
	}
	lu,e:=Login(name,group,string(pw))
	if e!=nil{
		pctx.Terminal.Println("\nsu: Login failed: "+e.Error())
		fc:=(pctx.Terminal.Environment)["_guest_su_auth_failed_count"].(int)
		fc++
		(pctx.Terminal.Environment)["_guest_su_auth_failed_count"]=fc
		if fc>=3 && (pctx.Terminal.Environment)["_guest_su_init"].(bool) == false{
			pctx.Terminal.Println("\nsu: Failed over 3 times.")
			terminal_dispose(pctx.Terminal)
		}
		return
	}
	(pctx.Terminal.Environment)["_guest_su_auth_failed_count"]=0
	(pctx.Terminal.Environment)["_guest_su_init"]=true
	pctx.Terminal.User=&lu
}
