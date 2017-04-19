# mmap_capnp test results

 - writing to mem-mapped file

The code for these tests can be found in [mmap_readwrite.go](mmap_readwrite.go)

The benchmark tests can be found in [mmap_readwrite_test.go](mmap_readwrite_test.go)

Writing capnp messages to memory mapped files, either individual messages or a single
list containing those messages. Messages have a minimal amount (12 bytes) of data


| Number of messages | List | Individual |
|---|---:|---:|
| 100 | 0.044 ms | 0.237 ms |
| 1.000 | 0.205 ms | 1.712 ms |
| 10.000 | 1.675 ms | 16.956 ms |
| 100.000 | 12.349 ms | 164.557 ms |
| 1.000.000 | 100.737 ms | 1452.672 ms |

According to this test, wrapping our messages in a list in another capnp object
achieves way better performance. Another very important observation is that writing
individual messages to the file requires a fixed size of the encoded message. This
makes it difficult to work with fields of variable length such as strings, lists...
Overal, the list method seems much better.

  - memory performance for holding messages in memory

The code for these tests can be found in [mem_usage.go](mem_usage.go)

There are 4 test cases:

  - store all messages in a capnp list.
  - store all messages in a capnp list, and encode the list
  - store all messages individually in a slice
  - store all messages encoded in a slice

This table contains the memory consumption of the respective data structures.

| number of messages | actual data in bytes per message| total data in bytes | capnp list | encoded capnp list | slice | slice with encoded data |
|:---|:---:|:---:|---:|---:|---:|---:|
| 1.000 | 12 B | 12 KB | 41 KB | 57 KB | 262 KB | 425 KB |
| 10.000 | 12 B | 120 KB | 172 KB | 336 KB | 2,8 MB | 1,6 MB |
| 100.000 | 12 B | 1,2 MB | 1,6 MB | 3,2 MB | 25,8 MB | 8,0 MB |
| 1.000.000 | 12 B| 12 MB | 16 MB | 32 MB | 259 MB | 81 MB |
| 1.000 | 22 B | 22 KB | 172 KB | 197 KB | 270 KB | 483 KB |
| 10.000 | 22 B | 220 KB | 3,1 MB | 3,4 MB | 3,0 MB | 2,25 MB |
| 100.000 | 22 B | 2,2 MB | 9,5 MB | 13MB | 25 MB | 12,0 MB |
| 1.000.000 | 22 B | 22 MB | 96 MB | 96 MB | 257 MB | 115 MB |
| 1.000 | 32 B | 32 KB | 270 KB | 311 KB | 303 KB | 500 KB |
| 10.000 | 32 B | 320 KB | 3,1 MB | 3,5 MB | 3,3 MB | 2,15 MB |
| 100.000 | 32 B | 3,2 MB | 7,9 MB | 16 MB | 28 MB | 15,4 MB |
| 1.000.000 | 32 B | 32 MB | 160 MB | 120 MB | 273 MB | 127 MB |
| 1.000 | 112 B | 112 KB | 2,2 MB | 2,3 MB | 598 Kb | 851 KB |
| 10.000 | 112 B | 1,12 MB | 4,7 MB | 4,8 MB | 4,2 MB | 3,25 MB |
| 100.000 | 112 B | 11,2 MB | 48 MB | 48 MB | 42 MB | 24,4 MB |
| 1.000.000 | 112 B | 112 MB | 360 MB | 480 MB * | 428 MB | 239 MB |
| 1.000 | 212 B | 212 KB | 3,0 MB | 3.3 MB | 876 KB | 1,33 MB |
| 10.000 | 212 B | 2.120 KB | 6,7 MB | 6.7 MB | 5,7 MB | 3,55 MB |
| 100.000 | 212 B | 21,2 MB | 90 MB | 67 MB | 50 MB | 35,4 MB |
| 1.000.000 | 212 B | 212MB | - | 896 MB * | 526 MB | 337 MB |
| 1.000 | 512 B | 512 KB | 2,6 MB | 3.1 MB | 1,8 MB | 2,53 MB |
| 10.000 | 512 B | 5.120 KB | 26 MB | 21 MB | 11,1 MB | 9,0 MB |
| 100.000 | 512 B | 51,2 MB | 156 MB | 208 MB | 80,8 MB | 91,4 MB |
| 1.000 | 1024 B | 1,0 MB | 3,0 MB | 4.1 MB | 3,4 MB | 2,53 MB |
| 10.000 | 1024 B | 10,0 MB | 41 MB | 31 MB | 20 MB | 15,25 MB |
| 100.000 | 1024 B | 100 MB | 310 MB | 413 MB * | 172 MB | 187,4 MB |

\* Decoding failed because the encoded message is too large

From these test results, we can see that the none-encoded list performance best if the
messages contain (almost) no data. As soon as we store increasing amounts of data in the
messages, our slice with the encoded data starts to perform better and better. All in all,
the slice method with encoded messages seems to have the best performance. Also note,
that the list method starts to suffer from decoding failures if we store too many
blocks. If the data used continues to grow, the none-encoded messages in the slice actually
seem to outperform the encoded messages. Whether to encode the messages or not should
probably be considered on a per case basis.

  - performance when loading messages from disk (with and without mmap)

The code for these tests can be found in [perf.go](perf.go)

The messages are stored in a single file. As per the findings of the first test,
the messages (actually instances of Tlog blocks) are stored in a capnp list
in a Tlog aggregation. The file is loaded, and we decode and check each block.
Results are averages over multiple runs.

| amount of blocks | file size in bytes | time to load with mmap | time to load without mmap |
| :--- | :---: | ---: | ---: |
| 1.000 | 16 KB | 0,30 ms | 0,26 ms |
| 10.000 | 160 KB | 2,94 ms | 2,81 ms |
| 100.000 | 1,6 MB | 7,29 ms | 6,77 ms |
| 1.000.000 | 16 MB | 33,63 ms | 31,52 ms |

Results show that mmap has a slightly negative impact on performance. This is likely
caused by the fact that we already store our messages in a single root object, which is
read and decoded in its entirety. Using mmap for reading seems to be redundant. 
