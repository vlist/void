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