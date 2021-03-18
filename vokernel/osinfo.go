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
		VoVersion:               "1.11.1 (2021.3.18)",
		GoVersion:               "go1.15.6 darwin/amd64",
		CurrentWorkingDirectory: getCwd(),
		SystemArch:              runtime.GOOS+"/"+runtime.GOARCH,
	}
}