# redis_crypto_perf

Test the time to hash, compress, encrypt, rehash, store in redis, load in redis,
decrypt and decompress binary objects.

The clients to connect to redis are imported from the [redis_perf](../redis_perf/redis)
test. The `redigo` client is the default, but both `radix` and `go-redis` are also
supported. For details on how to change the client and other parameters, run
`./redis_crypto_perf -h`.

## additional dependencies

These tests require a redis server. By default the application tries to connect on
`localhost:6379`. For redis v3 (the current stable release), the
[docker image](../redis_perf/Dockerfile) from the redis_perf is reused. For the
release candidate of redis v4, a [Dockerfile](Dockerfile) is provided in this directory.
The following commands produce a running instance listening on `localhost:6379`, with
a socket exposed in `/tmp/redis.sock`:

```bash
docker build -t customredis .
docker run --name redisv4 -p 6379:6379 -v /tmp:/tmp customredis
```

## workflow

***overview***

  1. Generate the data
  2. Store hash of data
  3. Compress the data
  4. Encrypt compressed data
  5. Store data in redis
  6. Reload data
  7. Decrypt data
  8. Decompress data

***generating the data***

To generate somewhat random data, the executed binary gets sliced up in the required
amount of sections of the requested length at the start of the application. If the
requested dataset is larger than the binary, sections are repeated. This produces
a somewhat random dataset that stays the same over different runs. Since the data
is not completely random, it is somewhat compressible and gives a relatively good
idea of the performance of the `gzip` and `lz4` compression ratios.


***hashing the data***

Hashing the objects occurs 2 times: immediately after they are generated, and after
they are encrypted. Hashing is done using the `md5` algorithm.

***compressing the data***

After the initial hash, the data gets compressed using either `gzip` or `lz4`. By
default the `gzip` algorithm is used, but `lz4` can be specified using a command line
flag. The compressed objects are stored separated from the original data, as the
original data is used at the end to verify the decrypted and decompressed objects.

***encrypting the data***

Once the objects are compressed, they are encrypted using the `AES-128` algorithm
running in `gcm` mode. The additional data used is hardcoded in the [encrypt.go](encryption/encrypt.go)
file. The encryption function automatically generates a random nonce for every object.
After the encryption, the nonce is prepended to the ciphertext (since it only needs to be
unique, not secret). When the ciphertext gets decrypted, the decrypt function first
reads the nonce from the ciphertext, and then continues to use the additional data
to reproduce the plaintext. The encryption key is the initial hash.

Including the nonce with ciphertext, as well as the authentication tag from the `gcm`
mode, causes the object size to increase when encrypted. The default nonce size is 12
bytes, while the size of the additional data can be changed if desired.

Once objects are encrypted, they overwrite the compressed objects in memory, and get
hashed again.

***storing the data***

All objects then get stored in redis. Every object gets stored in an `HSet` using the
hash of its encrypted form as the key and and the initial hash as the field. Although the
way the test data is generated causes some identical objects (therefore having an identical
initial hash and compressed form), the random nonce generation during the encryption
causes all objects to have different encrypted hashes which are used as keys. As a result,
all objects get stored in a seperate `HSet`.

After this, the objects are cleared from memory and reloaded. They then get decrypted and
decompressed again.
