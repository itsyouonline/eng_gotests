# redis_perf

A small application allowing users to test/benchmark different redis clients using
both tcp/ip and unix socket connections to redis.

## setup

The application expects a redis server to up and running before it is started. In this
directory, an example [redis.conf](redis.conf) file is provided, as well as a Dockerfile
that generates a redis image with said config file included. To create the image:

```bash
docker build -t customredis .
```
This will create an image with the name `customredis`

To then create the docker and set it up for our tests, run the command

```bash
docker run --name testredis -p 6379:6379 -v /tmp:/tmp customredis
```
A docker container will be created and started running redis, with the name testredis.
In the `/tmp` directory a Unix socket is exposed with the name `redis.sock`. The names
used here are just an example, and can be changed freely. To clear the storage, connect to
the container with the `redis-cli` and execute the `flushall` command:
```bash
docker exec -it testredis redis-cli
flushall
```

## usage

After running `go build`, an executable is created in this directory with the name
`redis_perf`. For a full list of all flags that can be used, run `./redis_perf -h`.
The most important flags are:

 - `-a, --object-amount`: the amount of objects to store in redis, default 1M.
 - `-s, --data-size`: the size per object in bytes, default 200 bytes.
 - `-t, --connectiontype`: the type of connection to redis, either `tcp` or `unix`, default `tcp`.
 - `-c, --connectionstring`: the connectionstring for redis, default `localhost:6379`.
 - `--client`: the underlying client to use, either `go-redis`, `redigo` or `radix`, default `go-redis`
 - `-p, --pipelength`: the amount of statements to store in a pipe before executing it.
 Setting this to 0 (or any negative number) disables the use of pipes. 0 by default,
 thus disabling pipes. 


## code

The code for the clients is found in the [redis](redis) directory. An interface `RedisClient`
is defined, and must be implemented by all clients. For every client, a struct is
defined that wraps the actual client. To implement an additional client, one should
define a new wrapper struct that holds the new client, and implement the `RedisClient`
interface. Finally, the constructor for the new client must be registered in the
`NewRedisClient` function in the [client.go](redis/client.go) file. The already
implemented clients can serve as an example on how to do this.

For every currently implemented client, there is also a `_test.go` file containing
a small benchmark using the go testing framework. These benchmarks simply create a new
client and measure how long it takes to store an item in a `HSet`. For the benchmark
tests in these files, the connectionstring cannot be specified. They expect a redis server
to be running on `localhost:6379` and a redis socket to be exposed at `/tmp/redis.sock`.
