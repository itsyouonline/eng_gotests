# Go test for capnp

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
