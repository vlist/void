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
		Version:               "1.12.3 (20A199)",
		Runtime_SystemArch:     runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH,
	}
}