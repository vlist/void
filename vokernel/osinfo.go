package vokernel

import (
	"runtime"
)

type OSInfo struct{
	Version               string
	Runtime_SystemArch    string
}
func GetOSInfo() OSInfo {
	return OSInfo{
		Version:               "1.12.1 (20A194)",
		Runtime_SystemArch:     runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH,
	}
}