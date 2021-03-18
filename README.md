# Void Shell
void shell is a CUSTOM shell service

## build voidshell
```shell
$ go clean
$ go build
```
## launch voidshell
```shell
$ ./void
```

## connect to voidshell
### using netcat
```shell
$ stty raw; nc -U ./voidsh
```
### using void socketterminal
```shell
$ ./sockterm ./voidsh
```
## configure voidshell
### socket file path
(in file .vsrc)
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
```diff
+ this will be highlighted in green
- this will be highlighted in red
```