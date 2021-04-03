package vokernel

import (
	"runtime"
)

type OSInfo struct{
	Version               string
	Runtime_SystemArch    string
}
//func getCwd() string{
//	wd,err:=os.Getwd()
//	if err!=nil{
//		return ""
//	}
//	return wd
//}
func GetOSInfo() OSInfo {
	return OSInfo{
		Version:               "1.12.01 (20A193d)",
		Runtime_SystemArch:     runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH,
	}
}