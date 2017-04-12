# tarantool test

## install tarantool

install on ubuntu 16.04
```
bash install_ubuntu.sh
```

## start tarantool

start tarantool without persistence
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


## build and start the test

```
# go build
# ./tarantool_perf -num=10000
test data:
	 number of messages:10000
	 data len = 200 bytes
	 time needed:1.582035538 seconds
```
