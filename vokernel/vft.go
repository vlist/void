package vokernel

import (
	"regexp"
	"strings"
)

/*
 * Void Format Text (VFT) 1.1.0
 * converts vft tag to terminal colors
 * only supports for forecolor and bold format
 * format: <vft {red|green|yellow|blue} {bold|}>formatting text</vft>,
   escape tag is <\vft> and <\/vft> (in string "<\\vft>","<\\/vft>")
 * example: black<vft red bold>red bold</vft>black<vft blue>blue</vft>black<\\vft green bold>shouldn't formatteded<\\/vft>
   output like: blackred boldblackblueblack<vft green bold>shouldn't formatted</vft>
*/
func Format(content string)string{
	reg:=regexp.MustCompile("<vft([\\S\\s]+?)>([\\S\\s]+?)</vft>")
	formatted:=reg.ReplaceAllStringFunc(content,func(i string)string{
		matchGroup:=reg.FindStringSubmatch(i)
		prop,inner:=matchGroup[1],matchGroup[2]
		var tagStart string="\033["
		if strings.Index(prop,"bold")!=-1{
			tagStart+="1;"
		}else{
			tagStart+="0;"
		}
		if strings.Index(prop,"red")!=-1{
			tagStart+="31m"
		}else if strings.Index(prop,"green")!=-1{
			tagStart+="32m"
		}else if strings.Index(prop,"yellow")!=-1{
			tagStart+="33m"
		}else if strings.Index(prop,"blue")!=-1{
			tagStart+="34m"
		}else{
			if strings.Index(prop,"bold")!=-1{
				tagStart="\033[1m"
			}else{
				tagStart="\033[0m"
			}
		}
		var tagEnd string="\033[0m"
		return tagStart+inner+tagEnd
	})
	formatted=strings.ReplaceAll(formatted,"<\\vft","<vft")
	formatted=strings.ReplaceAll(formatted,"<\\/vft","</vft")
	return formatted
}
