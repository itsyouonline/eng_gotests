# server

This package contains the [main.go](main.go), [tcpserver.go](tcpserver.go) and [storeserver.go](storeserver.go) files.
These 3 files make up (most of) the server source code. Once a `go build` command is executed in this directory, there will also be a `server` executable (assuming all dependencies have been downloaded correctly).
The server can optionally be started with a couple of command line flags:
- `-h, --help`: This will cause the server to print the application name, usage, version, commands and possible flags. The server then exits.
- `-d, --debug`: The server log level will be set to `debug`. This causes logs defined as debug to be printed to the terminal, which are ignored if this flag is not set. This has a serious impact on the benchmarks
- `-v, --version`: Like help, but only print the application name and version before exiting.
- `-b, --bind`: The port on which to bind the server. A ":" must be prepended to the port, else
an error will be thrown when starting the server

Note that all flags are optional. If the bind flag is not set, it defaults to the port specified in the constants.


## code overview


The `main.go` file contains some setup logic for the server.
The command flags get parsed, and some global constants are loaded in. Errors such as wrong flags get caught and cause the server to print the usage and exit.

the `storeserver.go` file contains the functions for the rpc calls:
create a new client, get and set

the `tcpserver.go` file contains the logic used to make a new server, and handle incomming connections


For a more in depth explanation of the code, see the comments in the respective files.
