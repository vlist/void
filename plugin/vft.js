/*
 * Void Format Text for Javascript (vft.js) 1.0
 * vft.go equivalence
 * converts vft tag to terminal colors
 * only supports for forecolor and bold format
 * format: <vft {red|green|yellow|blue} {bold|}>formatting text</vft>,
   escape tag is <\vft> and <\/vft> (in string "<\\vft>","<\\/vft>")
 * example: black<vft red bold>red bold</vft>black<vft blue>blue</vft>black<\\vft green bold>shouldn't formatteded<\\/vft>
   output like: blackred boldblackblueblack<vft green bold>shouldn't formatted</vft>
*/

function format(content){
    re="<vft([\\S\\s]+?)>([\\S\\s]+?)</vft>"
    reobj_glob=new RegExp(re,"g")
    reobj_tok=new RegExp(re)
    return content.replaceAll(reobj_glob,(sub)=>{
        grp=reobj_tok.exec(sub)
        prop=grp[1]; inner=grp[2]
        var tagStart="\033["
        if(prop.indexOf("bold")!==-1){
            tagStart+="1;"
        }else{
            tagStart+="0;"
        }
        if(prop.indexOf("red")!==-1){
            tagStart+="31m"
        }else if(prop.indexOf("green")!==-1){
            tagStart+="32m"
        }else if(prop.indexOf("yellow")!==-1){
            tagStart+="33m"
        }else if(prop.indexOf("blue")!==-1){
            tagStart+="34m"
        }else{
            if(prop.indexOf("bold")!==-1){
                tagStart="\033[1m"
            }else{
                tagStart="\033[0m"
            }
        }
        var tagEnd="\033[0m"
        return tagStart+inner+tagEnd
    })
        .replaceAll("<\\vft","<vft")
        .replaceAll("<\\/vft","</vft")
}

module.exports = {
    format: format
}