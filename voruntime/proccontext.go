package voruntime

import "void/vokernel"

type ProcContext struct{
	CommandName string
	Args        []string
	Type        string
	Terminal    *TerminalContext
	OS          vokernel.OSInfo
}
