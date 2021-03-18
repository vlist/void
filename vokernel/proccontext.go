package vokernel

type ProcContext struct{
	CommandName string
	Args []string
	Privileged bool
	Type string
	Shell *ShellContext
	OS OSInfo
}
