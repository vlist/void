package voruntime

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"void/vokernel"
)

var shmap=make(map[string]ListenerContext)

func internal_shutil(pctx *ProcContext) {
	rcsockid := "unix:" + vokernel.RC["socket"]
	if len(pctx.Args) == 0 {
		shutil_invalid_argument(pctx.Terminal)
		return
	}
	switch pctx.Args[0] {
	case "-o","--open":
		{
			var flag = ""
			if len(pctx.Args) <=1 {
				shutil_invalid_argument(pctx.Terminal)
				return
			}
			sockid := pctx.Args[1]
			if sockid == rcsockid {
				pctx.Terminal.Println("could not operate on default socket.")
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
					l, e = Startserver("unix", address,false)
				}
			case "tcp":
				{
					l, e = Startserver("tcp", address,false)
				}
			case "tls":
				{
					flag += "tls "
					println("starting server over TLS")
					l, e = Startserver_TLS("tcp", address)
				}
			default:
				{
					pctx.Terminal.Println("network " + network + " not supported.")
				}
			}
			if e != nil {
				pctx.Terminal.Println("opening shell on socket " + sockid + " failed.")
				log.Print(e)
				return
			}
			shmap[sockid] = ListenerContext{
				Listener: l,
				Flags:    flag,
			}
		}
	case "-k","--kill":
		{
			if len(pctx.Args) <=1 {
				shutil_invalid_argument(pctx.Terminal)
				return
			}
			sockid := pctx.Args[1]
			if sockid == rcsockid {
				pctx.Terminal.Println("could not operate on default socket.")
				return
			}
			if l, ok := shmap[sockid]; ok {
				e := (*l.Listener).Close()
				if e != nil {
					pctx.Terminal.Println("closing shell on socket " + sockid + " failed.")
					log.Print(e)
					return
				}
				delete(shmap, sockid)
			} else {
				pctx.Terminal.Println("closing shell on socket " + sockid + " failed: listener not found.")
			}
		}
	case "-l","--list":
		{
			pctx.Terminal.Println("opening socket shell: ")
			pctx.Terminal.Println(rcsockid + "\tdefault")
			for k, v := range shmap {
				pctx.Terminal.Println(k + "\t" + v.Flags)
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
	tctx.Println("invalid arguments.")
	usage := `usage [--options network:address]
options:
	-o,--open [tls|tcp|unix:address:port]: 
		create a new shell socket server
	-k,--kill [tls|tcp|unix:address:port]: 
		close specific socket server
	-l,--list: list all shell socket server
`
	tctx.Println(usage)
}