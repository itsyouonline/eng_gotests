# tarantool test

## install tarantool

install on ubuntu 16.04
```
bash install_ubuntu.sh
```

For other operating systems, see [the official download page](https://tarantool.org/download.html)

## start tarantool

start tarantool without persistence
```
bash clean.sh && tarantool start.lua
```

This starts an empty tarantool server on port 3302 (defined in [start.lua](start.lua))

Example output:
```
# bash clean.sh && tarantool start.lua
removed '00000000000000000000.xlog'
removed '00000000000000000000.vylog'
removed '00000000000000000000.snap'
2017-04-12 13:52:27.641 [3646] main/101/start.lua C> version 1.7.3-533-g84be303
2017-04-12 13:52:27.641 [3646] main/101/start.lua C> log level 5
2017-04-12 13:52:27.641 [3646] main/101/start.lua I> mapping 268435456 bytes for tuple arena...
2017-04-12 13:52:27.644 [3646] iproto/101/main I> binary: bound to 0.0.0.0:3301
2017-04-12 13:52:27.644 [3646] main/101/start.lua I> initializing an empty data directory
2017-04-12 13:52:27.648 [3646] snapshot/101/main I> saving snapshot `./00000000000000000000.snap.inprogress'
2017-04-12 13:52:27.649 [3646] snapshot/101/main I> done
2017-04-12 13:52:27.651 [3646] main/101/start.lua I> vinyl checkpoint done
2017-04-12 13:52:27.651 [3646] main/101/start.lua I> ready to accept requests
2017-04-12 13:52:27.652 [3646] main C> entering the event loop

```
After this, the tarantool server is ready to accept connections.

## build and start the test

```
# go build
# ./tarantool_perf -num=10000
```

Since we use the same space en id's every time we test, the tarantool server must
be restarted after every run as described in the [start tarantool](#start-tartantool) section.
The server and username/password to use can be supplied by command line flags:
	- `-user`: sets the tarantool username
	- `-passwd`: sets the tarantool password for the username
	- `-addr`: sets the server address, in the form of `IP_ADDRESS:PORT`

In addition, the following flags are accepted to change the test behaviour:
	- `num`: the amount of objects to store, defaults to 1M
	- `num-client`: the amount of concurrent goroutines to use, defaults to the amount of logical
	cpu cores available.

## Tarantool comparison with redis

When asked how to implement a HSet (as in redis), Tarantool developers responded
the following:

```
hello! tarantool is not key value.
it has spaces (aka table in rdbms) where you can store rows.
you can update row property as well (similar to hset)
```

They then offered the following example that came about as close as possible to
an HSet:

```lua
box.schema.create_space('xxx')
box.space.xxx:create_index('id', {parts = {1, 'string'}, type='HASH'})
box.space.xxx:put{'12346', 123456}
box.space.xxx:put{'12346', 146}
box.space.xxx:get{'12346'}
```

This is the approach we implemented in the test
