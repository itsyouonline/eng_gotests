# Results

For all tests, 200 bytes of random additional data are stored per block, next to the
36 bytes of data always present. See [the schema](capnp/tlog_schema.capnp) for the exact
block layout. In all tests, `AES` runs in 128bit mode (`AES-128`).

## Individual blocks

compress/decompress and encrypt/decrypt 10.000 blocks, total input size is 2.64 MB,
(264 bytes per block):

***compression and decompression***

| algorithm | compressed size | ratio | time taken (compress) | input bytes/second (compress) | time taken (decompress) | input bytes/second (decompress) |
| --- | --- | --- | --- | --- | --- | --- |
| lz4  | 2,83 MB | +7,19% | 7,96 s | 331 KB/s | 21,24 s | 133 KB/s |
| gzip | 2,84 MB | +7,76% | 309,17 ms | 8,54 MB/s | 60,05 ms | 47,29 MB/s |
| zlib | 2,72 MB | +3.03% | 343,20 ms | 7,69 MB/s | 64,70 ms | 42,04 MB/s |

***encryption and decryption***

| algorithm | mode of operation | time taken (encrypt) | bytes/second (encrypt) | time taken (decrypt) | bytes/second (decrypt) |
| --- | --- | --- | --- | --- | --- |
| AES | cbc | 10,50 ms | 251,43 MB/s | 7,48 ms | 363,64 MB/s |
| AES | cfb | 9,41 ms | 280,55 MB/s | 7,55 ms | 349,67 MB/s |
| AES | ctr | 7,87 ms | 335,45 MB/s | 5,60 ms | 471,43 MB/s |
| AES | ofb | 5,83 ms | 452,83 MB/s | 8,00 ms | 330,00 MB/s |
| AES | gcm | 4,17 ms | 633,09 MB/s | 7,09 ms | 394,92 MB/s |
| 3DES | cbc | 563,65 ms | 4,68 MB/s | 766,73 ms | 3,55 MB/s |
| 3DES | cfb | 945,83 ms | 2,79 MB/s | 483,99 ms | 5,56 MB/s |
| 3DES | ctr | 449,33 ms | 5,88 MB/s | 450,60 ms | 5,86 MB/s |
| 3DES | ofb | 451,84 ms | 5,84 MB/s | 450,87 ms | 6,21 MB/s|
| 3DES | gcm | - | - | - | - |
| Twofish | cbc | 37,46 ms | 70,48 MB/s | 37,03 ms | 73,45 MB/s |
| Twofish | cfb | 36,17 ms | 72,99 MB/s | 45,55 ms | 57,96 Mb/s |
| Twofish | ctr | 38,55 ms | 68,48 MB/s | 35,20 ms | 75,00 MB/s |
| Twofish | ofb | 38,21 ms | 69,09 MB/s | 36,26 ms | 72,80 MB/s |
| Twofish | gcm | 55,69 ms | 47,40 MB/s | 54,13 ms | 51,73 MB/s |
| Blowfish | cbc | 31,09 ms | 84,91 MB/s | 36,75 ms | 74,01 MB/s |
| Blowfish | cfb | 30,40 ms | 86,84 MB/s | 36,75 ms | 71,84 MB/s |
| Blowfish | ctr | 29,43 ms | 89,70 MB/s | 26,77 ms | 98,62 MB/s |
| Blowfish | ofb | 28,22 ms | 93,55 MB/s | 26,48 ms | 99,70 MB/s|
| Blowfish | gcm | - | - | - | - |

## individual blocks, Concatenated

compress/decompress and encrypt/decrypt 1.000.000 blocks, total input size is 264 MB,
(264 bytes per block):

***compression and decompression***

| algorithm | compressed size | ratio | time taken (compress) | input bytes/second (compress) | time taken (decompress) | input bytes/second (decompress) |
| --- | --- | --- | --- | --- | --- | --- |
| lz4  | 237,59 MB | -10,01% | 1,87 s | 141,18 MB/s | 488,81 ms | 486,05 MB/s |
| gzip | 224,91 MB | -14,81% | 8,44 s | 31,28 MB/s | 3,29 s | 68,36 MB/s |
| zlib | 224,92 MB | -14,80% | 8,53 s | 30,95 MB/s | 3,33 s | 67,54 MB/s |

***encryption and decryption***

| algorithm | mode of operation | time taken (encrypt) | bytes/second (encrypt) | time taken (decrypt) | bytes/second (decrypt) |
| --- | --- | --- | --- | --- | --- |
| AES | cbc | 545,86 ms | 483,64 MB/s | 467,18 ms | 565,09 MB/s |
| AES | cfb | 596,11 ms | 442,87 MB/s | 615,17 ms | 429,15 MB/s |
| AES | ctr | 408,73 ms | 645,90 MB/s | 428,58 ms | 615,17 MB/s |
| AES | ofb | 367,86 ms | 717,66 Mb/s | 327,24 ms | 806,75 MB/s |
| AES | gcm | 112,56 ms | 2,35 GB/s | 108,83 ms | 2,43 GB/s |
| 3DES | cbc | 46,29 s | 5,70 MB/s | 45,66 s | 5,78 MB/s |
| 3DES | cfb | 45,09 s | 5,85 MB/s | 44,63 s | 5,92 MB/s |
| 3DES | ctr | 44,17 s | 5,98 MB/s | 44,05 s | 5,99 MB/s |
| 3DES | ofb | 44,42 s | 5,94 MB/s | 43,43 s | 6,09 MB/s |
| 3DES | gcm | - | - | - | - |
| Twofish | cbc | 3,43 s | 76,97 MB/s | 3,06 s | 86,27 MB/s |
| Twofish | cfb | 3,08 s | 85,71 MB/s | 3,06 s | 86,27 MB/s |
| Twofish | ctr | 3,13 s | 84,35 MB/s | 3,12 s | 84,62 MB/s |
| Twofish | ofb | 3,15 s | 83,81 MB/s | 3,44 s | 76,74 MB/s |
| Twofish | gcm | 4,54 s | 58,15 MB/s | 4,37 s | 60,41 Mb/s |
| Blowfish | cbc | 2,52 s | 104,76 MB/s | 2,65 s | 99,62 MB/s |
| Blowfish | cfb | 2,65 s | 99,62 MB/s | 2,65 s | 99,62 MB/s |
| Blowfish | ctr | 2,24 s | 117,86 MB/s | 2,24 s | 117,86 MB/s |
| Blowfish | ofb | 2,30 s | 114,78 MB/s | 2,29 s | 115,28 MB/s |
| Blowfish | gcm | - | - | - | - |

## blocks stored in list

compress/decompress and encrypt/decrypt 1.000.000 blocks stored in a single list,
total input size is 248MB:

***compression and decompression***

| algorithm | compressed size | ratio | time taken (compress) | input bytes/second (compress) | time taken (decompress) | input bytes/second (decompress) |
| --- | --- | --- | --- | --- | --- | --- |
| lz4  | 235,81 MB | -4,92% | 561,11 ms | 441,98 MB/s | 412,28 ms | 571,96 MB/s |
| gzip | 222,55 MB | -10,27% | 7,14 s | 34,73 MB/s | 557,92 ms | 398,89 MB/s |
| zlib | 222,55 MB | -10,26% | 7,29 s | 34,02 MB/s | 661,62 ms | 336,37 MB/s |

***encryption and decryption***

| algorithm | mode of operation | time taken (encrypt) | bytes/second (encrypt) | time taken (decrypt) | bytes/second (decrypt) |
| --- | --- | --- | --- | --- | --- |
| AES | cbc | 468,01 ms | 529,90 MB/s | 362,49 ms | 684,16 MB/s |
| AES | cfb | 485,98 ms | 510,31 MB/s | 468,42 ms | 529,44 MB/s |
| AES | ctr | 344,77 ms | 719,32 MB/s | 347,01 ms | 714,68 MB/s |
| AES | ofb | 325,34 ms | 762,28 Mb/s | 303,00 ms | 818,48 MB/s |
| AES | gcm | 102,13 ms | 2,43 GB/s | 87,95 ms | 2,82 GB/s |
| 3DES | cbc | 41,50 s | 5,98 MB/s | 41,32 s | 6,00 MB/s |
| 3DES | cfb | 41,68 s | 5,95 MB/s | 41,62 s | 5,96 MB/s |
| 3DES | ctr | 41,30 s | 6,00 MB/s | 41,34 s | 5,99 MB/s |
| 3DES | ofb | 41,31 s | 6,00 MB/s | 41,23 s | 6,02 MB/s |
| 3DES | gcm | - | - | - | - |
| Twofish | cbc | 3,23 s | 76,78 MB/s | 2,78 s | 89,21 MB/s |
| Twofish | cfb | 2,88 s | 86,11 MB/s | 2,86 s | 86,71 MB/s |
| Twofish | ctr | 2,92 s | 84,93 MB/s | 2,92 s | 84,93 MB/s |
| Twofish | ofb | 2,98 s | 83,22 MB/s | 2,96 s | 83,78 MB/s |
| Twofish | gcm | 4,18 s | 59,33 MB/s | 4,34 s | 57,14 MB/s |
| Blowfish | cbc | 2,37 s | 104,64 MB/s | 2,29 s | 108,30 MB/s |
| Blowfish | cfb | 2,45 s | 101,22 MB/s | 2,47 s | 100,40 MB/s |
| Blowfish | ctr | 2,09 s | 118,66 MB/s | 2,10 s | 118,10 MB/s |
| Blowfish | ofb | 2,15 s | 115,35 MB/s | 2,14 s | 115,89 MB/s |
| Blowfish | gcm | - | - | - | - |

## conclusion

Compressing individual capnp messages is very counterproductive. If the message can't
be compressed, it actually grows in size (due to the compression headers). Also,
a lot of time is spend, while (trying to) compress.

That being said, `lz4` clearly outperforms `gzip` and `zlib` in terms of speed for large
continuous messages. As a downside, it does offer a reduced compression ratio. For small
individual messages, `gzip` and `zlib` performs better, especially when decompressing, with `gzip` appearing slightly better than `zlib`. Overall, whether to apply compression or not should probably be decided on a case by case basis.
Due to the lack of a real dataset of sufficient size, the additional data stored here
is generated randomly, thus it isn't compressible by definition. Should this extra data be
structured (such as natural text), compressing might be worthwhile. In this sense,
these results are in line with [the official capnp documentation on the subject](https://capnproto.org/encoding.html#compression).


In case of encryption, the `AES-128` algorithm clearly and always outperforms the others.
As for the modes of operation: `gcm` is the fastest for both encrypting and decrypting,
while also offering the added security of authenticated encryption. It does require
the entire message to be known upfront and encrypts it at once. `ofb`, `ctr` and `cfb`
are stream modes, but unauthenticated. `cbc` is the only block mode, and is also
unauthenticated. Input for `cbc` MUST be padded, and thus the produced output will be
slightly larger than the input (16 bytes / encrypted object at worst). Likewise, in `gcm`
mode, the implementation provided by the go standard library will append the tag to the
encrypted output by default, thereby growing every encrypted message with 12 bytes.

For encrypting and decrypting, `AES-128` in `gcm` mode should be preferred as it is both the
fastest and authenticated.
