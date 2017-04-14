# client

This package contains the [main.go](main.go) and [storeclient.go](storeclient.go) files.
These 2 files make up (most of) the client source code. Once a `go build` command is executed in this directory, there will also be a `client` executable (assuming all dependencies have been downloaded correctly).
The client can optionally be started with a couple of command line flags:
- `-h, --help`: This will cause the client to print the application name, usage, version, commands and possible flags. The client then exits.
- `-d, --debug`: The client log level will be set to `debug`. This causes logs defined as debug to be printed to the terminal, which are ignored if this flag is not set.
- `-v, --version`: Like help, but only print the application name and version before exiting.
- `-s, --server`: The value following this flag is set as the remote capnp rpc server address to be used by this client in the benchmark. The port must be set to the same value as the server

Note that all flags are optional. If the server flag is not set, it defaults to localhost and the port specified in the constants.


## code overview


The `main.go` file contains some setup logic for the client.
The command flags get parsed, and some global constants are loaded in. Errors such as wrong flags get caught and cause the client to print the usage and exit.

the `storeclient.go` file contains the main benchmark logic:
open and close a connection to the server, and execute the rpc benchmarks.

For a more in depth explanation of the code, see the comments in the respective files.
