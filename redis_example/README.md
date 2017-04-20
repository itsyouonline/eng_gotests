# redis_example

Example code showing how to store capnp messages in a redis hset and reading them again.
To compile the example, run `go build`. A separate redis instance is required. Example setup
of a docker running redis:

```bash
docker run --name redis -p 6379:6379 redis
```
This starts a docker container with name `redis` that listens on `localhost:6379`.
Note that data in this redis instance is not persistent.

## code overview

in [main.go](main.go) we do mostly setup (parsing the command flags, setting the logger).
At the end of the file, a new redis client instance is created. This client is then
passed to the `storeAndReadCapnpInHset` function in the [redis.go file](redis.go),
allong with the amount of messages we want to store. This function creates the required
amount of messages, encodes them and stores them in the redis hset. When this is done,
some of the blocks are being retrieved again from the hset and decoded. The sequence
of the decoded blocks is checked to make sure it is as we expect it to be. If there is
a discrepancy, an error is logged and the program exits. If debug output is enabled,
a log message will also be printed for every block that has passed the check.

## flags

  - `-d, --debug`: enables additional debug output
  - `-c, --connectionstring`: redis connection string, defaults to `localhost:6379`
  - `-l, --data-length`: the amount of bytes of extra data to store in the tlog blocks.
  Default is 0.
  - `-a, --message-amount`: the amount of message to store in this run, default is 1.000
  - `-h, --help`: show the command line help output
  - `-v, --version`: print the version information
