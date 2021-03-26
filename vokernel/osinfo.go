package vokernel

import (
	"os"
	"runtime"
)

type OSInfo struct{
	VoVersion               string
	GoVersion               string
	CurrentWorkingDirectory string
	SystemArch              string
}
func getCwd() string{
	wd,err:=os.Getwd()
	if err!=nil{
		return ""
	}
	return wd
}
func GetOSInfo() OSInfo {
	return OSInfo{
		VoVersion:               "1.11.4 (20A146)",
		GoVersion:               runtime.Version()+" "+runtime.GOARCH,
		CurrentWorkingDirectory: getCwd(),
		SystemArch:              runtime.GOOS+"/"+runtime.GOARCH,
	}
}