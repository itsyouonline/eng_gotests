# Go test for capnp

This directory contains the go source files for the capnp benchmark. The following structure is in place:
- client:
  - [main.go](client/main.go)
  - [storeclient.go](client/storeclient.go)
  - (after building) the client executable
- constants:
  - [common.go](constants/common.go)
- server:
  - [main.go](server/main.go)
  - [storageserver.go](server/storageserver.go)
  - [tcpserver.go](server/tcpserver.go)
  - (after building) the server executable
- store:
  - [specification.capnp.go](store/specification.schema.capnp)


- [generate.go](generate.go)
- [readme.md](readme.md)


## Regeneration of the capnp generated files

First make sure you have the capnp tools installed:

**Mac osx**

```
brew install capnp
```

**Linux**

```
apt-get install capnproto
```

And the go capnp generation library:
```
go get -u zombiezen.com/go/capnproto2/...
```

### Regenerate

```
go generate
```
