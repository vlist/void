# voidshell
voidshell is a CUSTOM shell service
![avatar](void.png)
Current version: 1.11.1 (2021.3.18).<br/>
Author: jlywxy (ms2692848699@outlook.com)<br/><br/>
IMPORTANT: this program now don't support Windows. [see reason](#miscellaneous)<br/>

## build voidshell
```shell
$ go clean
$ go build
```
## launch voidshell
```shell
$ ./void
```
To launch in background, type
```shell
$ screen -R voidsh
$ ./void
```
...or make a voidshell.service file to let systemctl manage voidshell service on linux.

## connect to voidshell
### using netcat
```shell
$ stty raw; nc -U ./voidsh
```
### using socketterminal
See https://github.com/jlywxy/socketterminal
```shell
$ ./socketterminal ./voidsh
```
voidshell listen to one default unix socket connections only(when launching).<br/>
To change that default unix socket file path, modify [configuration files](#shutil).</br>
For multiple kind of terminal connecting concurrently, use that default socket to configure, see builtin command [shutil](#shutil).

## configure voidshell
Configuraton file: .vsrc
### socket file path
<span id="conffile.sock"></span>
```json
{
  "socket": "./socketfile"
}
```
### password for internal command "sudo"
"password_encrypted" should be sha256 encrypted.
```json
{
  "password_encrypted": "sha256(password)"
}
```
### plugin root directory
```json
{
  "plugin_root": "path"
}
```
Plugin files should be located in path/root/ firectory.

## plugin development
Plugins are node.js script file,
located in plugin/root directory.
### create a plugin
Create a javascript file located in plugin/root directory.<br/>
Plugin template is shown below:
```javascript
/*init code,do not modify*/
var ctx={};module.exports={ init: (_ctx)=>{ctx=_ctx}, run: main }

function main(){
    ctx.print("Hello void.")
    ctx.exit()
}
```
Plugin entrypoint is function main.
### using plugin
Simply type plugin name and arguments in voidshell.
```shell
void:>plugins-name
Hello void.
```
### calling convention

Data comes with `ctx`: 
* plugin name and arguments: <br/>
   `ctx.args=[plugin_name, plugin_arg1, plugin_arg2, ...]`
* plugin root directory: <br/>
   `ctx.pwd="./plugins/root"`
* input: <br/>
   `ctx.input(prompt,callback_func)`
* output: <br/>
  output a string:
   `ctx.print(content)`<br/>
  output a VFT formatted string:
   `ctx.printf(content)`<br/>
* void format text transformer: <br/>
   `ctx.format(text)`
* exit the plugin: <br/>
  `ctx.exit()`
  
### void format text(VFT)
* Converts vft tag to VT100 terminal colors.
* Only supports forecolor and bold format.
* Format:`<vft {red|green|yellow|blue} {bold|}>formatting text</vft>`,<br/>
  escape tag is `"<\vft>"` and `"<\/vft>"`
* Example: `"black<vft red bold>red bold</vft>black<vft blue>blue</vft>black<\vft green bold>shouldn't formatteded<\/vft>"`,<br/>
  output should be:
  black<span style="color: red; font-weight: bold">red bold</span>black<span style="color: blue">blue</span>black&lt;vft green bold&gt;shouldn't formatted&lt;/vft&gt;<br/>
* Different implements in vft.go and vft.js are equivalent.
  
## builtin commands
### info
Displays os info and shell/terminal info.
```shell
void:>info

                    _      __  __           
     _   __ ____   (_) ___/ / _\ \          
    | | / // __ \ / // __  / (_)\ \         
    | |/ // /_/ // // /_/ / _   / / ______    
    |___/ \____//_/ \____/ (_) /_/ /_____/  
     void:>void --everything

Void System 1.1
    Golang Version: go1.15.6 darwin/amd64
    Current Working Directory: *
    System Arch: darwin/amd64
Process Context(pctx):
    Command Name: info
    Arguments: []
    Shell Context(sctx): 
        Terminal Name: *
        Privileged: false
```
### exec
Run bash commands in voidshell.
```shell
void:>exec ls
README.md main.go   plugins    void.png  voidsh    vokernel  voruntime voshell
```
### shutil
Manage socket servers.<br/>
```shell
void:>shutil
usage [--options network address]
  --open: create a new shell socket server
  --kill: close specific socket server
  --list: list all shell socket server
```
network: "tcp" or "unix"(unix socket).<br/>
address: "ip:port" for "tcp", or socket filename for "unix".<br/><br/>
Examples:<br/>
Opening new socket servers:
```shell
void:>shutil --open tcp:127.0.0.1:9001
void:>shutil --open unix:/tmp/vssock1
```
Close a socket server:
```shell
void:>shutil --kill unix:/tmp/vssock1
````
List all opened socket server:
```shell
void:>shutil --list
opening socket shell: 
unix:./voidsh (default)
tcp:127.0.0.1:9001
unix:/tmp/vssock1
```
The default socket neither could be reopened nor be killed.
### exit
Simply exit the shell(close current terminal only, won't shut down service).
```shell
void:>exit
```
<br/>
Other commands are now reserved for further development,<br/>
see in voruntime/internal.go

## miscellaneous
* voidshell and socketterminal(https://github.com/jlywxy/socketterminal) use the protocol of VT100 terminal.
* voidshell now DO NOT support windows, because voidshell use unix socket to listen and run initializing commands, while Windows don't support unix socket.

## update log
1.11.2 (20A0319)  *Newest Alpha
* added a configuration option to set plugin root directory, see [voidshell configuration](#configure-voidshell)
* modified plugin calling conventions, see [plugin development](#plugin-development)
* internal code modifications:
    * change `voruntime.Exec` to `voruntime.BashExec`, then `voruntime.Exec` directly spawn process while `voruntime.BashExec` use /bin/bash to eval bash commands
    * added keys `ctx.printf` and `ctx.pwd` to plugin_init.js
    * made `voruntime.Process` able to pass plugin root in `voruntime.RC` to plugin_init.js 

1.11.1 (2021.3.18 20:00)
* added builtin command "shutil" to manage server sockets. 
now voidshell could accept multiple kinds of terminal connections concurrently
* internal code modifications:
    * made `net.Listen(network,path)` reusable in socketshell.go, added map `shmap` to save server listeners, added function of shutil in map `voruntime.internal`

1.11 (2021.3.18 17:37) bugfix
* removed dependency of xterm-resize in module "exec.go"; add "resize.go" for equivalence
* add usage of builtin command "exec" in "README.md"
* add some miscellaneous for terminal protocol in "README.md"
    * removed bash commands `resize>>/dev/null` in `voruntime.Exec`, then added function `voruntime.Getsize` in resize.go to get terminal rows and cols.

1.0 (2021.3.18)
repo initialization
