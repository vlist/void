package voruntime

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var shmap=make(map[string]ListenerContext)

func internal_shutil(pctx *ProcContext) {
	rcsockid := "unix:" + RC["socket"]

	if len(pctx.Args) == 0 {
		shutil_invalid_argument(pctx.Terminal)
		return
	}
	switch pctx.Args[0] {
	case "--open":
		{
			var flag = ""
			if len(pctx.Args) < 2 {
				shutil_invalid_argument(pctx.Terminal)
				return
			}
			sockid := pctx.Args[1]
			if sockid == rcsockid {
				pctx.Terminal.Output("could not operate on default socket\n")
				return
			}
			na := strings.Split(sockid, ":")
			network := na[0]
			address := strings.Join(na[1:], ":")

			var l *net.Listener
			var e error
			switch network {
			case "unix":
				{
					os.RemoveAll(address)
					dir := filepath.Dir(address)
					err := os.MkdirAll(dir, 0770)
					if err != nil {
						println("checking directory failed: " + dir)
						e = err
					}
					l, e = Startserver("unix", address)
				}
			case "tcp":
				{
					l, e = Startserver("tcp", address)
				}
			case "tls":
				{
					flag += "tls "
					println("starting server over TLS")
					l, e = Startserver_TLS("tcp", address)
				}
			default:
				{
					pctx.Terminal.Output("network " + network + " not supported\n")
				}
			}
			if e != nil {
				pctx.Terminal.Output("opening shell on socket " + sockid + " failed\n")
				log.Print(e)
				return
			}
			shmap[sockid] = ListenerContext{
				Listener: l,
				Flags:    flag,
			}
		}
	case "--kill":
		{
			if len(pctx.Args) < 2 {
				shutil_invalid_argument(pctx.Terminal)
				return
			}
			sockid := pctx.Args[1]
			if sockid == rcsockid {
				pctx.Terminal.Output("could not operate on default socket\n")
				return
			}
			if l, ok := shmap[sockid]; ok {
				e := (*l.Listener).Close()
				if e != nil {
					pctx.Terminal.Output("closing shell on socket " + sockid + " failed\n")
					log.Print(e)
					return
				}
				delete(shmap, sockid)
			} else {
				pctx.Terminal.Output("closing shell on socket " + sockid + " failed: listener not found\n")
			}
		}
	case "--list":
		{
			pctx.Terminal.Output("opening socket shell: \n")
			pctx.Terminal.Output(rcsockid + "\tdefault\n")
			for k, v := range shmap {
				pctx.Terminal.Output(k + "\t" + v.Flags + "\n")
			}
		}
	default:
		{
			shutil_invalid_argument(pctx.Terminal)
			return
		}
	}
}

func shutil_invalid_argument(tctx *TerminalContext){
	tctx.Output("invalid argument")
	usage := `usage [--options network:address] [--tls]
options:
	--open: create a new shell socket server
	--kill: close specific socket server
	--list: list all shell socket server
--tls:
	serve over TLS
`
	tctx.Output(usage)
}