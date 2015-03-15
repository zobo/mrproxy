go-mrproxy
=======

A memcache to redis proxy/protocol translator written in golang

#### Building

```
go get -u github.com/zobo/mrproxy
go get -u github.com/garyburd/redigo
cd $GOPATH/src/github.com/zobo/mrproxy
```

###### 64bit:
- Linux:   `GOOS=linux GOARCH=amd64 go build -o mrproxy.bin -v mrproxy/main.go`
- OSX:     `GOOS=darwin GOARCH=amd64 go build -o mrproxy.bin -v mrproxy/main.g`
- Windows: `GOOS=windows GOARCH=amd64 go build -o mrproxy.bin -v mrproxy/main.go`

###### 32bit:
- Linux:   `GOOS=linux GOARCH=386 go build -o mrproxy.bin -v mrproxy/main.go`
- OSX:     `GOOS=darwin GOARCH=386 go build -o mrproxy.bin -v mrproxy/main.go`
- Windows: `GOOS=windows GOARCH=386 go build -o mrproxy.bin -v mrproxy/main.go`
