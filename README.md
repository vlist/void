# voidshell
voidshell is a CUSTOM shell service

## build voidshell
```shell
$ go clean
$ go build
```
## launch voidshell
```shell
$ ./void
```
to launch in background, type
```shell
$ screen -R voidsh
$ ./void
```
...or make a voidshell.service file to let systemctl manage voidshell service on linux.

## connect to voidshell
voidshell is now listening unix socket connections only.
### using netcat
```shell
$ stty raw; nc -U ./voidsh
```
### using void socketterminal
(now under development, see https://github.com/jlywxy/sockterm)
```shell
$ ./sockterm ./voidsh
```
## configure voidshell
configuraton file: .vsrc
### socket file path

```json
{
  "socket": "./socketfile"
}
```
### password for internal command "sudo"
"password_encrypted" should be sha256 encrypted
```json
{
  "password_encrypted": "sha256(password)"
}
```

## plugin development
plugin for voidshell is node.js script file,
located in plugin/root directory
### calling convention

1. plugin args: <br/>
   `ctx.args[plugin_name, plugin_arg1, plugin_arg2, ...]`
2. input: <br/>
   `ctx.input(prompt,callback_func)`
3. output: <br/>
   `ctx.print(content)`
4. exit plugin: <br/>
   `ctx.exit()`
5. transform void format text: <br/>
   `ctx.format(text)`
   
### void format text(VFT)
* converts vft tag to terminal colors
* only supports forecolor and bold format
* format:`<vft {red|green|yellow|blue} {bold|}>formatting text</vft>`,<br/>
  escape tag is `<\vft>` and `<\/vft>` (in string `"<\\vft>"`,`"<\\/vft>"`)
* example: `black<vft red bold>red bold</vft>black<vft blue>blue</vft>black<\\vft green bold>shouldn't formatteded<\\/vft>`<br/>
  output like:
  black<span style="color: red; font-weight: bold">red bold</span>black<span style="color: blue">blue</span>black&lt;vft green bold&gt;shouldn't formatted&lt;/vft&gt;
  
## builtin commands
### info
displays os info and shell/terminal info
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

### exit
simply exit the shell(close terminal only, won't shut down service)
```shell
void:>exit
```
<br/>
other commands are now reserved for further development.<br/>
see in voruntime/internal.go: internal