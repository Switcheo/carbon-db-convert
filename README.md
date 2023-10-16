# db-convert

## What it does
Converts LevelDB database files into RocksDB. Adapted from level-to-rocks to work with Carbon, and uses cometbft-db plus grocksdb for compatibility with rocksdb 7.10.2.

## Prerequisites
* Go >= 1.18.x
* Existing LevelDB database in Carbon
* RocksDB >= 7.10.2

## How to build
```sh
$ make build
```

## How to use
1. make sure carbon is initialised with goleveldb as the backend, and has `.db` files in `.carbon/data`
2. run the binary, specifying custom in/out directories using the flags if needed, and rocksdb `.db` files will be created in the output directory

    ```sh
    $ ./db-convert -dbDir="database file dir" -outDir="output file dir"
    ```
3. copy and replace the goleveldb `.db` files with the generated ones, then change the `db_backend` option in `.carbon/config/config.toml` to rocksdb
4. restart carbon (ensure that it is built with rocksdbflags), and verify that the chain is continuing from the previous block height