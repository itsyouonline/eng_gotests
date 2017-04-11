package main

//go:generate capnp compile -I$GOPATH/src/zombiezen.com/go/capnproto2/std -ogo:store ../specification.capnp
// For some stupid reason capnp does not put it in the right output folder
//go:generate mv specification.capnp.go store/
