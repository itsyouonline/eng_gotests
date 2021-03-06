# gotests

This repo contains the code for the go tests. Code is grouped in different directories
for different tests, e.g. all the mmap code is in one subdirectory, and the redis code is in another.

## installing

First create the directory structure:
```bash
mkdir -p $GOPATH/src/github.com/itsyouonline
```

Go looks for packages we import in very specific places, therefore it is important
to immediately put our source code in the right directory to avoid issues later on.

Now clone our repo with the sources:
```bash
cd $GOPATH/src/docs.greenitglobe.com/itsyouonline
git clone git@github.com:itsyouonline/eng_gotests.git
cd eng_gotests
```

As a last step, we need to download the dependencies:
```bash
go get ./...
```

When this is done, we have all the required code to start building

Note that none code dependencies (such as the tarantool server, etc...) install
instructions can be found in the README.md files in the respective sub directories

## building

To build the code for a group of tests, go into the directory and, from a terminal,
execute the `go build` command. This creates an executable in that directory. Information
about specific tests/executables can be found in the README.md file in their respective
directories. Example:

```bash
cd mmap_capnp
go build
```

This will create the `mmap_capnp` executable in the `mmap_capnp` directory.
Example:
![build example](buildexample.png)

An exception to this is the `capnp_benchmarks` directory. More details can be found
there in the [capnp_benchmarks' README.md for go](capnp_benchmarks/go/README.md)

### code documentation

Although all code is documented, code in the  `capnp_benchmark` subdirectory has excessive
documentation to introduce people new to golang. If at any point something should not be clear,
[the official go documentation](https://golang.org/doc/effective_go.html) will
probably have the answer
