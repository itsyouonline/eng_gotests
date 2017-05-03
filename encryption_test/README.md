# encryption_test

This directory contains all of the code for the encryption/decryption and
compression/decompression tests.

## what is tested

For compressing and decompressing, the following algorithms have been considered:

  - lz4
  - zlib
  - gzip

The implementations from the go runtime have been used for zlib and gzip. For lz4,
the [pierrec/lz4](https://github.com/pierrec/lz4) package is used.

For encryption and decryption, the following have been considered:

  - AES-128
  - 3DES
  - Blowfish
  - Twofish
  - RC4 (rejected due to known vulnerabilities)

Additionally, these ciphers can run in different modes, being:

  - CBC
  - CFB
  - CTR
  - OFB
  - GCM*

All of these encryption algorithms and modes of operation are provided by the go
runtime.

*REMARKS:*

  - AES runs in 128 bit mode due to the 16 byte key that is being used.
  - 3DES needs a 24 instead of 16 byte key.
  - GCM is the only authenticated mode of operation, which is recommended for strong security.
  - GCM requires a 16byte block size, thus it is not compatible with blowfish and 3DES.

All tests are run on capnp messages containing tlog blocks, specified in the
[capnp file](capnp/tlog_schema.capnp). For each test, the time it takes for every
compression and decompression algorithm is measured. After decompression, the data
is verified to be the same as the input for the compression. Then, the encryption
and decryption is tested with every possible algorithm in every possible mode of
operation. Likewise, data integrity is checked after decryption. Encryption is checked
against the (original) uncompressed data. There are 3 test scenarios:

  - an amount of tlog blocks, each stored in a separate message, compressed and encrypted separately
  - an amount of tlog blocks, each stored in a separate message, then encoded in byte form and
  concatenated. The concatenated array is then compressed and encoded as a whole. This requires a fixed size
  for encoded messages to be restored again to a capnp format.
  - an amount of tlog blocks, stored in a list in a single message, compressed and encrypted as a whole.

## running the tests

After running the `go build` command, this directory will contain an executable called
`encryption_test`. Invoking this application will run all the tests and print (intermediate)
results to the terminal. Additionally some flags can be supplied:

  - `--debug, -d`: Enable debug logging
  - `--help, -h`: Show help
  - `--version, -v`: Print the version
  - `--size, -s`: Set the amount of random data to store in a block (next to the fixed size fields
    that are always present)
  - `--amount, -a`: Amount of blocks to use in the test. For the test with individual blocks, this
  amount is divided by 100 
