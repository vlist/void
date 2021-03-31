package voruntime

import "strings"

func PermissionFilter(commandname string,permission string) (bool,string){
	if permission == ""{
		return true,"All commands access are granted of this user group."
	}
	permissionToken:=strings.Split(permission,",")
	tmpPerm:=false
	tmpReason:="Permission restrictions of command "+commandname+" not found."
	for _,v:=range permissionToken{
		if v=="-"{
			tmpPerm=false
			tmpReason="Permission denied for all commands of this user group."
		}else{
			if v=="-"+commandname {
				tmpPerm=false
				tmpReason="Permission denied for command "+commandname+" of this user group."
			}else if v==commandname{
				tmpPerm=true
				tmpReason="Command "+commandname+" is granted of this user group."
			}else{
				continue
			}
		}
	}
	return tmpPerm,tmpReason
}
