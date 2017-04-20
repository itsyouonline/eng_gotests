# mmap_capnp test results

The capnp schema file used for these tests can be found [here](tlog_schema.capnp)

## writing to mem-mapped file

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
Overall, the list method seems much better.

## memory performance for holding messages in memory

The code for these tests can be found in [mem_usage.go](mem_usage.go)

There are 4 test cases:

  - store all messages in a capnp list.
  - store all messages in a capnp list, and encode the list
  - store all messages in a capnp list. Judge the required buffer size and provide an empty buffer with that capacity
  - store all messages in a capnp list, and encode the list. Judge the required buffer size and provide an empty buffer with that capacity
  - store all messages individually in a slice
  - store all messages encoded in a slice

This table contains the memory consumption of the respective data structures:

First column: the amount of messages (actually tlog blocks) we keep in memory.

Second column: the amount of actual data per message/block.

Third column: the total amount of actual data we store
Remainder: the size of the respective data structure to store said data.

| number of messages | actual data in bytes per message| total data in bytes | capnp list | encoded capnp list |  capnp list with buf | encoded capnp list with buf | slice with messages | slice with encoded data |
|:---|:---:|:---:|---:|---:|---:|---:|---:|---:|---:|
| 1.000 | 12 B | 12 KB | 41 KB | 57 KB | 33 KB | 49 KB | 262 KB | 425 KB |
| 10.000 | 12 B | 120 KB | 172 KB | 336 KB | 328 KB | 492 KB | 2,8 MB | 1,6 MB |
| 100.000 | 12 B | 1,2 MB | 1,6 MB | 3,2 MB | 3,2 MB | 4,7 MB | 25,8 MB | 8,0 MB |
| 1.000.000 | 12 B| 12 MB | 16 MB | 32 MB | 32 MB | 48 MB | 259 MB | 81 MB |
| 1.000 | 22 B | 22 KB | 172 KB | 197 KB | 57KB | 90 KB | 270 KB | 483 KB |
| 10.000 | 22 B | 220 KB | 3,1 MB | 3,4 MB | 562 KB | 909 KB | 3,0 MB | 2,25 MB |
| 100.000 | 22 B | 2,2 MB | 9,5 MB | 13MB | 5,7 MB | 8,9 MB | 25 MB | 12,0 MB |
| 1.000.000 | 22 B | 22 MB | 96 MB | 96 MB | 57 MB | 89 MB | 257 MB | 115 MB |
| 1.000 | 32 B | 32 KB | 270 KB | 311 KB | 82 KB | 123 KB | 303 KB | 500 KB |
| 10.000 | 32 B | 320 KB | 3,1 MB | 3,5 MB | 844 KB | 1,2 MB | 3,3 MB | 2,15 MB |
| 100.000 | 32 B | 3,2 MB | 7,9 MB | 16 MB | 8,3 MB | 12 MB | 28 MB | 15,4 MB |
| 1.000.000 | 32 B | 32 MB | 160 MB | 120 MB | 58 MB | 98 MB | 273 MB | 127 MB |
| 1.000 | 112 B | 112 KB | 2,2 MB | 2,3 MB | 352 KB | 483 KB | 598 KB | 851 KB |
| 10.000 | 112 B | 1,12 MB | 4,7 MB | 4,8 MB | 3,6 MB | 4,8 MB | 4,2 MB | 3,25 MB |
| 100.000 | 112 B | 11,2 MB | 48 MB | 48 MB | 19 MB | 31 MB | 42 MB | 24,4 MB |
| 1.000.000 | 112 B | 112 MB | 360 MB | 480 MB * | 191 MB | 311 MB * | 428 MB | 239 MB |
| 1.000 | 212 B | 212 KB | 3,0 MB | 3.3 MB | 655 KB | 885 KB | 876 KB | 1,33 MB |
| 10.000 | 212 B | 2,12 MB | 6,7 MB | 6,7 MB | 2,6 MB | 4,8 MB | 5,7 MB | 3,55 MB |
| 100.000 | 212 B | 21,2 MB | 90 MB | 67 MB | 34 MB | 59 MB | 50 MB | 35,4 MB |
| 1.000.000 | 212 B | 212MB | 672 MB | 896 MB * | 362 MB | 586 MB *| 526 MB | 337 MB |
| 1.000 | 512 B | 512 KB | 2,6 MB | 3.1 MB | 1,6 MB | 2,1 MB | 1,8 MB | 2,53 MB |
| 10.000 | 512 B | 5,12 MB | 26 MB | 21 MB | 8,7 MB | 13 MB | 11,1 MB | 9,0 MB |
| 100.000 | 512 B | 51,2 MB | 156 MB | 208 MB | 89 MB | 141 MB | 80,8 MB | 91,4 MB |
| 1.000 | 1024 B | 1,0 MB | 3,0 MB | 4.1 MB | 3,1 MB | 4,1 MB | 3,4 MB | 2,53 MB |
| 10.000 | 1024 B | 10,0 MB | 41 MB | 31 MB | 18 MB | 27 MB | 20 MB | 15,25 MB |
| 100.000 | 1024 B | 100 MB | 310 MB | 413 MB * | 179 MB | 282 MB * | 172 MB | 187,4 MB |

\* Decoding failed because the encoded message is too large

First, the list and encoded list method without the buffer (the library handles allocating
a new buffer), has terrible cpu performance for writing a new list as soon as we encode
of data. That being said, providing a buffer with an estimated size also seems to increase
memory performance significantly. If we use a capnp list, encoding it consitently
increases the memory consumption, so in this case it seems best to just hold the
capnp message containing the Tlog aggregation with said list in memory none-encoded.
The slice on the other hand, seems to perform better when the messages are encoded and
stored, on most occasions. Note that the slice approach does not perform well when the
actual amount of stored data per message is low. As soon as we approach ~100 bytes of data
per message, efficiency increases to the point where it seems similar to our list.
The benefit of our slice is that messages are kept separate, whereas they are combined
into 1 message with our list approach. As soon as the message with the list gets above
about 250 MB in size, the library refuses to decode it back. The slice could potentially store
gigabytes worth of data. Whichever approach is better will probably depend on the
requirements of the application, either estimating the size of the objects to be
generated and creating a list with a buffer of said size, or storing the (possibly
encoded) messages separately in a slice.

## performance when loading messages from disk (with and without mmap)

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
read and decoded in its entirety. Using mmap for the solo purpose of reading seems to be redundant.
