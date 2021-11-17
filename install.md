# Install

## From release

*> Linux:*
```shell
curl -lO -L https://github.com/ariary/fileless-xec/releases/latest/download/fileless-xec
```

*> Windows:*
```shell
curl -lO -L https://github.com/ariary/fileless-xec/releases/latest/download/fileless-xec_windows.exe
```

## From source

Clone the repo and download the dependencies locally:
```    
git clone https://github.com/ariary/fileless-xec.git
cd fileless-xec
make before.build
```

To build the fileless-xec for linux :

     build.fileless-xec
    
To build the fileless-xec for windows :

     windows.build.fileless-xec

## With `go` command

Make sure `$GOPATH` is in your `$PATH` before

Install `fileless-xec`

```shell
go install github.com/ariary/fileless-xec/cmd/fileless-xec
```
