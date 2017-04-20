# Tarantool performance results

For these tests, tarantool was running on localhost:3302

Every object we store has a size of 200 bytes.

It is worth noting that there is currently only 1 tarantool client. This client is
developed by the Tarantool team.

Tests have been run with a varying amount of objects and goroutines. The goroutines
each open their own connection to the tarantool sever. Goroutines are created sequentially
and start storing data as soon as they are created. Timing starts *before* the routines are created.

| Amount of objects | Amount of goroutines used | Time taken in seconds | Items stored per second | Items stored per second per goroutine |
| --- | :---: | ---: | ---: | ---: |
| 1.000 | 1 | 0,285 s | 3.508 | 3.508 |
| 10.000 | 1 | 1,846 s | 5.417 | 5.417 |
| 100.000 | 1 | 8,521 s | 11.735 | 11.735 |
| 1.000.000 | 1 | 33,001 s | 30.302 | 30.302 |
| 1.000 | 2 | 0,153 s | 6.535 | 3.268 |
| 10.000 | 2 | 0,355 s | 28.169 | 14.084 |
| 100.000 | 2 | 2,441 s | 40.966 | 20.483 |
| 1.000.000 | 2 | 19,716 s | 50.720 | 25.360 |
| 1.000 | 4 | 0,125 s | 8.000 | 2.000 |
| 10.000 | 4 | 0,253 s | 39.525 | 9.881 |
| 100.000 | 4 | 1,395 s | 71.684 | 17.921 |
| 1.000.000 | 4 | 15,305 s | 65.338 | 1.6334 |
| 1.000 | 8 | 0,099 s | 10.101 | 1.262 |
| 10.000 | 8 | 0,204 s | 49.019 | 6.127 |
| 100.000 | 8 | 1,033 s | 96.805 | 12.100
| 1.000.000 | 8 | 9,351 s | 106.940 | 13.367 |

**Important:** As per the tarantool developers, Tarantool `is not key value`. This example
has been implemented to mimic a `redis hset` as closely as possible. For more info, see
the [README.md](README.md) file.

Per the results, storing more items is more efficient (more items/second, since the overhead of
opening the connection is proportionally smaller compared to the time spend actually sending items).
Using multiple goroutines is less efficient per goroutine (because we start sending items as soon as
the routine is created, therefore the first routine has already send some items when the second one starts
etc...), however the global efficiency increases (as expected) since more items get send.
For applications that gradually store objects, using a single goroutine will be best,
but if a larger amount of messages must be stored at once, multiple goroutines should be used.
(though this increases cpu cost).
