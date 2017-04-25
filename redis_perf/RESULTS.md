# redis_perf test results

Test setup:

  - Redis instance running in a docker on localhost:
    - Container port `6379` forwarded
    - Mounted `/tmp` to expose the unix socket
    - [redis config file](redis.conf)
  - All tests store 1M objects
  - Every object has a size of 200 bytes

3 different clients have been tested:

  - [go-redis](github.com/go-redis/redis)
  - [radix.v2](github.com/mediocregopher/radix.v2)
  - [redigo](github.com/garyburd/redigo)

For every client, both `Tcp/ip` and `unix socket` connections to redis are tested.
Also for every client and every connection type, tests are run with and without
pipelining. After every test, the database is flushed.

## results

pipe size `-` indicates that no pipe was used.

| client | connection type | pipe size | total time | objects per second |
| :---: | :---: | ---: | ---: | ---: |
| go-redis | tcp/ip | - | 41,627 s | 24.022 obj/s |
| go-redis | tcp/ip | 1 | 41,445 s | 24.128 obj/s |
| go-redis | tcp/ip | 5 | 9,592 s | 104.253 obj/s |
| go-redis | tcp/ip | 10 | 6,016s | 166.223 obj/s |
| go-redis | unix socket | - | 13,843 s | 72.239 obj/s |
| go-redis | unix socket | 1 | 13,639 s | 73.319 obj/s |
| go-redis | unix socket | 5 | 5,879 s | 170.097 obj/s |
| go-redis | unix socket | 10 | 3,671 s | 272.405 obj/s |
| redigo | tcp/ip | - | 32,289 s | 30.970 obj/s |
| redigo | tcp/ip | 1 | 32,756 s | 30.529 obj/s |
| redigo | tcp/ip | 5 |  8,129 s | 123.016 obj/s |
| redigo | tcp/ip | 10 | 5,772 s | 173.250 obj/s |
| redigo | unix socket | - | 9,961 s | 100.391 obj/s |
| redigo | unix socket | 1 | 10,060 s | 99.403 obj/s |
| redigo | unix socket | 5 | 3,760 s |  265.957 obj/s |
| redigo | unix socket | 10 |  2,862 s | 349.406 obj/s |
| radix.v2 | tcp/ip | - | 32,302 s | 30.958 obj/s |
| radix.v2 | tcp/ip | 1 | 32,493 s | 30.776 obj/s |
| radix.v2 | tcp/ip | 5 | 12,739 s | 78.499 obj/s |
| radix.v2 | tcp/ip | 10 | 9,388 s | 106.519 obj/s |
| radix.v2 | unix socket | - | 10,685 s | 93.589 obj/s |
| radix.v2 | unix socket | 1 | 11,111 s | 90.001 obj/s |
| radix.v2 | unix socket | 5 | 4,804 s | 208.160 obj/s |
| radix.v2 | unix socket | 10 | 3,664 s | 272.926 obj/s |

For these timing results, only the time taken to store the data and check for an
error while doing so is considered, stored items weren't read from redis.

The results clearly show that unix sockets vastly outperform tcp connections.
This is especially true when pipes aren't used, as the rtt cost has to be payed
for every single object. Considering this, unix sockets should always be preferred
where applicable.

The better client to use is dependent on whether or not pipelining is used, and if so,
how much commands are used per pipeline. In case no pipeline is used, the redigo and
radix.v2 clients are roughly tied concerning speed, with go-redis lagging behind.

In case pipelining is used, redigo outperforms the other tested clients. In tests with
pipelining, the responses are ignored, we only look for errors. Since radix requires us
to load in every response from the pipeline and check it for an error, it looses a lot of
performance compared to redigo and go-redis, which get any error from anywhere in the pipe
when reading the first response or when executing the pipe. They then ignore the actual
command response, since we know it succeeded anyway as there was no error. In Radix, a
possible error is part of the response for a command, therefore we must read every response
in the pipe and check if it has an error.
