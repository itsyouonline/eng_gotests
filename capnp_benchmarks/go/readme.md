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

## Usage

The code consists of 2 separate parts, which must be compiled independently. There is both a server and client. Compiling is as simple as going into each of the server and client directory and executing the `go build` command. For the server, from this directory:
```bash
cd server
go build
```
Likewise for the client, from this directory:
```bash
cd client
go build
```

This will create a runnable `server` and `client` executable in the respective directory. To get started with the benchmarks, first start the server by going to the server directory, and start the server:
```bash
cd server
./server
```
Now that the server is running, it is possible to start one or more clients to run the benchmarks:
```bash
cd client
./client
```

When a client gets started, it will automatically connect to the server and start the benchmarks. At the end of each benchmark, the time taken gets printed in the terminal. Once the benchmarks are done, the client exits.

For more details on both the client and the server, see the `readme.md` files in their respective directories.



## Regeneration of the capnp generated files
**Important**: The code used is schema dependent. If the [schema](../specification.capnp) is changed, and the generated files are regenerated, the code _**wil**_ break.

First make sure you have the capnp tools installed:

**Mac osx**

```
brew install capnp
```

**Linux**

```
apt-get install capnproto
```

And the go capnp generation library (this should already be installed as it is required to compile the code):
```
go get -u zombiezen.com/go/capnproto2/...
```

### Regenerate

```
go generate
```
