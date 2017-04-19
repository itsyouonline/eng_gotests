# Capnp nim benchmark

## Install required packages

Install capnp compiler
```
nimble install capnp
ln -s ~/.nimble/bin/capnpc ~/.nimble/bin/capnpc-nim
```

Install reactor.nim and collections.nim

```
nimble install collections
nimble install reactor
```

## compile the spec
```
capnp compile -onim specification.capnp > specification.nim
```

There is a modification in the spec. `struct Object` changed to `struct Product`,
because `Object` is Nim keyword

TODO:
it seems we can still use `Object`

## Run server
```
nim c -r server.nim
```

## Run client

```
nim c -r client.nim
```
