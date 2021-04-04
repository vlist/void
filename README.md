# voidshell
voidshell is a CUSTOM shell service
![avatar](void.png)
Current version: 1.12.1 (20A194). [See update log](#update-log).<br/>
Author / Contributors: <a href="https://github.com/jlywxy">jlywxy</a>, <a href="https://github.com/vlist">vlist</a>.
<br/><br/>
This program now don't support Windows. [see reason](#windows-no)<br/>

## prerequisite
python3 and python3-dev should be installed.<br/>
package for Ubuntu: python3-dev<br/>
for CentOS: python3-devel<br/>
for Mac OS: (no need to install -dev)

## build voidshell
```shell
$ ./build.sh
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
## develop voidshell by Goland
If using Goland:
try Run directly. 
If `pkg-config: exit status 1`: run 
```
$ ./cgo-cpython-tool.sh


...
...
------------------------------
(if needed) Set in .bash_profile or Goland Run/Debug Configuration -> Environment:
PKG_CONFIG_PATH=/

Put in plugin.go file:
#cgo pkg-config: python3
------------------------------
```
Set `#cgo pkg-config` in voruntime/plugin.go and PKG_CONFIG_PATH shown in the result to Goland Run/Debug Configurations -> Environment

## connect to voidshell
### using voidterminal
See https://github.com/jlywxy/voidterminal
```shell
$ ./voidterminal unix:/tmp/vssock1
```
### using netcat
```shell
$ stty raw; nc -U /tmp/vssock1; stty -raw
```
voidshell open one default unix socket (/tmp/vssock1) server only when launching.<br/>
To change that default socket path, modify [configuration files](#shutil).</br>
When connected, use builtin command [shutil](#shutil) to open other unix,tcp or tls shell servers.

## configure voidshell
Configuraton file: vsrc.json
### socket file path
<span id="conffile.sock"></span>
```json
{
  "socket": "./tmp/vssock1"
}
```
Plugin files should be located in path/root/ firectory.
### server TLS certificate
```json
    {
      "tls_config_pem": "cert/server.pem",
      "tls_config_key": "cert/server.key"
    }
```
IMPORTANT: cert files cert/server.* in this repository are self-signed and should NEVER be used in production mode. <br/>
## Configure users and permission
User and Permission Filter scheme is helpful to avoid attackers to harm the host via voidshell.
User data are stored in users.json.<br/>
Every user has its group, which is easy to control Permission.<br/>
Permission Filter Syntax:<br/>
`""` for all granted,<br/>
`"-"` for all denied,</br>
`"-shadow"` for denying access to `shadow`<br/>
`"-,shutil"` for denying all, then allow `shutil`<br/>
Permission Layer: internal(I),exec(E),plugin(P)<br/>
Note: Exec Permission should be `""`(all granted) or `"-"`(all denied) ONLY, <br/>
because `exec` parameters are directly passed to /bin/bash, void interpreter cannot analyse what processes are being spawned. (user may use `;`,`&`,`&&`,`|`,`screen`,etc to spawn multiple process in single line; user also could spawn another shell to avoid the permission filtering.).<br/>
<br/>
Permissions set for single user is a implement of group permission.However, leaving `""` in single user permission setting won't overwrite group permission to "all granted".<br/>
Password_Encrypted is sha256 hashed.<br/>
DO NOT use "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" as password_encrypted because it means empty string: sha256("")<br/>
Example:<br/>
```json
{
  "user_groups":
  {
    "admin": {
      "admin": {
        "password_encrypted": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "permission": {"internal":"","exec":"","plugins":""}
      }
    },
    "guest": {
      "guest": {
        "password_encrypted": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "permission": {"internal":"-who","exec":"","plugins":""}
      }
    }
  },
  "group_permission":
  {
    "admin": {"internal":"","exec":"","plugins":""},
    "guest": {"internal":"-,info,su,who,clear,exit,_stop_repl","exec":"-","plugins":"-"}
  }
}
```
Note: _stop_repl should be added to internal permissions, which is a dependency of internal command `shadow`.<br/>
Users connected to default socket shell will login to admin:admin automatically.<br/>
When terminal transmission is not secured (raw tcp), [su](#su) is disabled to prevent password leakage. To allow `su` in insecured transmission, modify configuration in `vsrc.json`:
```json
{
  "allow_su_via_insecure_transmission": "true"
}
```
### configure plugin root directory
```json
{
  "plugin_root": "./plugins"
}
```

## plugin development
Plugins are python3 script file,
located in plugin/root directory.
### create a plugin
Create a python file located in plugin/root directory.<br/>
Plugin template is shown below:
```
#DO NOT modify "init"
def init(pctx):
    global ctx
    ctx=pctx

def main(args):
    ctx.print("Hello void.\n")
```
Plugin entrypoint is function main.
### using plugin
Simply type plugin name and arguments in voidshell.
```shell
void:> plugin
Hello void.
```
### developing plugin
`ctx` provides most informations and functions to be used for plugin.
* `ctx.print(content)`<br/>
Output `content` to current terminal stdout
  
* `ctx.println(content)`<br/>
Output `content` and Line Feed to current terminal stdout

* `ctx.printf(content)`<br/>
Output `content` wrapped by VFT tags to current terminal stdout
  
* `ctx.input(prompt)`<br/>
Scan input after output `prompt`
  
* `ctx.tctx`<br/>
Readonly dict representing current terminal context.<br/>
Type "pinfo" to check terminal context passed to plugin.
```shell
void:admin# pinfo
voidshell Plugin Extension 1.0
Terminal Context: 
{'tctxid': '902ba40c-0f09-ff26-4987-1d34c24b2a05', 
'tctx': {'RawConnection': {}, 'StdinReader': {}, 
'StdoutWriter': {}, 'StdinWriterSwitch': {'Destination': {}}, 
'Secured': True, 'Delim': 13, 'ShellName': 'unix:/tmp/vssock1', 
'TerminalID': '902ba40c-0f09-ff26-4987-1d34c24b2a05', 
'User': {'Name': 'admin', 'Group': 'admin', 'Permission': 
[',', ',', ',']}, 'Environment': {'_guest_su_auth_failed_count': 0,
 '_guest_su_init': True}}, 'root': './plugins/root/',
  'self_root': './plugins/root/'}

[void]: 'pinfo' exit 0.
```
  
* `ctx.root`,`ctx.self_root`<br/>
"root" represents plugin root directory (set in vsrc.json config file),<br/>
"self_root" represents root directory of "this" script.<br/>
  Note: when plugin form is plugin_name.py, self_root is its parent directory; (plugin name: path/plugin, root=path/, self_root=path/)<br/>
  When plugin form is plugin_name/\__init\__.py, self_root is the parent directory of \__init\__.py (plugin name: path/plugin, root=path/, self_root=path/plugin)<br/>
  
### void format text(VFT)
* Converts vft tag to VT100 terminal codes.
* Only supports forecolor and bold format.
* Format:`<vft {red|green|yellow|blue} {bold|}>formatting text</vft>`,<br/>
  escape tag is `"<\vft>"` and `"<\/vft>"`
* Example: `"black<vft red bold>red bold</vft>black<vft blue>blue</vft>black<\vft green bold>shouldn't formatteded<\/vft>"`,<br/>
  output should be:
  black<span style="color: red; font-weight: bold">red bold</span>black<span style="color: blue">blue</span>black&lt;vft green bold&gt;shouldn't formatted&lt;/vft&gt;<br/>
  
## builtin commands
### info
Display os info,terminal info,user info.
```shell
void:admin# info

                    _      __  __           
     _   __ ____   (_) ___/ / _\ \          
    | | / // __ \ / // __  / (_)\ \         
    | |/ // /_/ // // /_/ / _   / / ______    
    |___/ \____//_/ \____/ (_) /_/ /_____/  
     void:> void --everything

voidshell 1.12.1 (20A194)
└─ Runtime/System Arch: go1.16.2 darwin/amd64
Process Context(pctx):
├─ Command Name: info
├─ Arguments: []
└─ Terminal Context(tctx):
   ├─ Shell Interface: unix:/tmp/vssock1
   ├─ Terminal ID: 11bec1e0-683e-af5e-8e6a-67b4d9bfba4b
   ├─ Transmission Secured: true
   └─ User Context(uctx):
      ├─ User Identifier: admin:admin
      └─ Permissions: I(),E(),P(): all granted

```
Arguments could also be applied to control `info` output.
```shell
void:admin# info --nologo --noctx
voidshell 1.12.1 (20A194)
└─ Runtime/System Arch: go1.16.2 darwin/amd64
```
### exec
Run bash commands in voidshell.
```shell
void:> exec ls
README.md           plugin_def.h
__pycache__         plugins
build.sh            users.json
cert                void
cgo-cpython-tool.sh void.png
go.mod              vokernel
go.sum              voruntime
main.go             vsrc.json
```
### shutil
Command version 1.2.<br/>
Manage socket servers.<br/>
```shell
void:> shutil
usage [--options network:address]
options:
	-o,--open [tls|tcp|unix:address:port]: 
		create a new shell socket server
	-k,--kill [tls|tcp|unix:address:port]: 
		close specific socket server
	-l,--list: list all shell socket server
```
network: "tcp" or "unix"(unix socket).<br/>
address: "ip:port" for "tcp", or socket filename for "unix".<br/><br/>
To configure TLS certificate, see [TLS Certificate Configuration](#server-tls-certificate)
Examples:<br/>
Opening new socket servers:
```shell
void:> shutil --open unix:/tmp/vssock2
void:> shutil --open tls:127.0.0.1:9001
void:> shutil -o tcp:127.0.0.1:9001
```
Close a socket server:
```shell
void:> shutil --kill unix:/tmp/vssock2
void:> shutil -k tls:127.0.0.1:9001
````
List all opened socket server:
```shell
void:> shutil --list
opening socket shell: 
unix:/tmp/vssock1       default
tls:127.0.0.1:9000  tls
tcp:127.0.0.1:9001
unix:/tmp/vssock2
```
The default socket neither could be reopened nor be killed.
### shadow
Command version 1.1.<br/>
Project current terminal output to another terminal.The terminal that be projected is called "shadow terminal".<br/>
terminal id can be looked up in [info](#info).<br/>
```shell
void:> shadow
usage [--commands] [terminal id]
commands:
	-p,--project [terminal id]
		project current terminal session to specific terminal
	-d,--detach
		detach shadow terminal
```
Examples:
in main terminal:
```shell
void:> shadow --project 3f35f171-eaf6-c9c0-3021-60049d96d4ba
void:> shadow -d
close existing shadow projector: 3f35f171-eaf6-c9c0-3021-60049d96d4ba
```
in shadow terminal:
```shell
shadow connecting to: 436605b1-06de-33ae-7a03-23aa03e361fa
--------SHADOW BEGINS--------

void:> _stop_repl
void:> shadow --detach
--------SHADOW ENDS--------
shadow disconnecting from: 436605b1-06de-33ae-7a03-23aa03e361fa
```
```_stop_repl``` is a reserved internal command which is used to pause REPL of the current terminal.<br/>
When terminal became a shadow terminal, its stdin will write to itself other than write to the main terminal, which means the shadow terminal could only see but not respond to anything.<br/>
### su
Switch user. To configure users, see [configure users and permission](#configure-users-and-permission).
```shell
void:> su admin:jlywxy
su: Enter password for admin:jlywxy: 
void:jlywxy#
```
### exit
Simply exit the terminal.
```shell
void:>exit
disconnected
```
<br/>
Other commands are now reserved for further development,<br/>
see in voruntime/internal.go

## miscellaneous
* voidshell and socketterminal(https://github.com/jlywxy/socketterminal) use the protocol of VT100 terminal.
### Windows, NO
* voidshell now DO NOT support windows, because
    1. voidshell use unix socket to listen configuration terminal session when launching, 
    while Windows don't support unix socket.
    2. voidshell build tool now can only run on *nix system.
    3. voidshell running in WSL have not been tested.
    
## update log
1.12.1 (20A194) *Newest Alpha
* fixed issues #7.
* made plugin experience more stable.
* changed syntax to configure permission.[configure permission](#configure-users-and-permission)
* removed all legacy code of node.js plugins.
* now could login to voidshell directly from stdio, but it will mess with voidshell log output.

1.12.01 (20A193d) *Newest Alpha-dev
* added internal command `su` to switch user: [su](#su),added internal `who` to get user info.
* plugin now base on cgo python3, which is experimental.
* a new build tool is used. (./build.sh)
* input supports history track (use up and down arrow key), history files are located in .voidsh_history/terminalid
* internal code modification:
  * added a internal command __cast_admin_uuid in case admin password forgotten.
  
1.11.4 (20A146) 
* changed command syntax of [shutil](#shutil) and [shadow](#shadow).
* fixed bugs of command shadow.
* internal code modification:
    * made [shutil](#shutil) and [shadow](#shadow) case match long and short flags(eg: ```--detach```to```-d```)
    * added function ```StartREPL``` and ```StopREPL``` to ```TerminalContext``` to make terminal REPL loop controllable.
* version name rule changed: 1.11.4 (20A146) represents to: <br/>
  major version 1; framework version 1, functions version 1; bugfix version 4, dev version 0(not a dev version); <br/>
  build id: year 2020; state Alpha; build date(hex): 0x146 (hex(0326)). <br/>
    
1.11.31 (20A0325d) 
* issue fix
    * automatically mkdir of directories in given unix socket file path which are not exist.
* added internal command "shadow", see [shadow](#shadow).
* added identifier for every terminal session(Terminal Name).
* changed syntax of "shutil" to open a TLS server. See [shutil](#shutil).
* internal code modifications:
    * changed `ShellContext` to `TerminalContext` in all related code, which is making more sense.
    * moved some code (terminalcontext.go, proccontext.go, etc) from `vokernel` layer to `voruntime` layer.<br/>
    * changed `TerminalContext` structure.
    
1.11.3 (20A0323) 
* added TLS support
* internal code modifications:
    * added `Startserver_TLS` func to `socketshell.go`
    * modified internal command `shutil` to make available to open TLS socket server.</br>
    
1.11.2 (20A0319)  
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
