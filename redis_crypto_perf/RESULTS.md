# redis_crypto_perf test results

All tests are run with 1M objects of 200 bytes each. These objects are generated
from the binary, therefore if the amount and size of the objects is the same in
multiple runs, the dataset will be identical. This way sufficiently random data
is generated while still being relatively compressible.

The redis instances ran in a docker, with the [redis.conf config file](redis.conf)
found in this directory. The different redis backends used are redis v3.2 and
redis v4.0-rc3. A [dockerfile for redis v4](Dockerfile) is also found in this directory.

All tests are run with the redigo client, while connecting to the redis via a unix
socket. The performance analysis of different clients and sockets vs tcp-ip was
conducted earlier and can be found [here](../redis_perf/RESULTS.md).

Per the [results in encryption_test](../encryption_test/RESULTS.md), the chosen
cryptographic algorithm is AES-128 in gcm mode. The additional data has a length
of 19 bytes. A random nonce is generated for every object during the encryption.

The workflow is described in more detail in [the README file](README.md).

## results

| compression algorithm | redis version | time to hash original objects | compression time | compressed size | encryption time | encrypted size | time to hash encrypted objects | storage time | redis memory used to store objects | loading time | decryption time | decompression time |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| gzip | v3 | 453 ms | 284,35 s | 157,59 MB (78,8%) | 1,58 s | 185,59 MB | 478 ms | 11,85 s | 485,24 MB (261,46%) | 11,89 s | 1,01 s | 17,90 s |
| lz4 | v3 | 500 ms | 906,49 s | 190,82 MB (95.41%) | 1,57 s | 218,82 MB | 505 ms | 11,47 s | 516,37 MB (235,98%) | 10,49 s | 938 ms | 30 m 41,37 s |
| gzip | v4 | 463 ms | 242,54 s | 157,59 MB (78,8%) | 1,50 s | 185,59 MB | 447 ms | 9,99 s | 453,32 MB (244,26%) | 9,05 s | 921 ms | 14,37 s |
| lz4 | v4 | 460 ms | 900,18 s | 190,82 MB (95,41%) | 1,48 s | 218,82 MB | 475 ms | 9,98 s | 484,37 MB (221.35%) | 9,29 s | 945 ms | 24 m 51.86 s |


First of all, it is clear that the golang implementation of `lz4` is very  bad compared
to the `gzip` implementation in the standard library. Regardless of the implementation used,
compression takes up most of the time spend.

The upcoming v4 release of redis does have improved memory performance, but overall
redis takes up more than double the amount of memory that we actually use. Memory
performance is increased when larger objects are stored.  
