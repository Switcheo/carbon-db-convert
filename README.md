# carbon-db-convert

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
    $ ./carbon-db-convert -dbDir="database file dir" -outDir="output file dir"
    ```
3. copy and replace the goleveldb `.db` files with the generated ones, then change the `db_backend` option in `.carbon/config/config.toml` to rocksdb
4. restart carbon (ensure that it is built with rocksdb flags), and verify that the chain is continuing from the previous block height

## Convert existing carbon chain from leveldb to rocksdb

1. Install rocksdb from source

    Ensure [dependencies](https://github.com/Switcheo/rocksdb/blob/v7.10.2-patched/INSTALL.md#dependencies) are met.

    ```bash
    # install dependencies for Ubuntu:
    sudo apt-get install build-essential cmake libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev libzstd-dev -y
    ```

    ```bash
    # install rocksdb
    cd ~
    git clone https://github.com/Switcheo/rocksdb.git
    cd rocksdb
    make shared_lib
    sudo make install-shared
    sudo ldconfig
    ```

2. Install carbon-db-convert
    ```bash
    git clone https://github.com/Switcheo/carbon-db-convert.git
    cd carbon-db-convert
    make build
    ```

3. Ensure you have enough space for conversion. Below are an estimation of additional space required for /rocksdata convert size.

    ```bash
    # before
    4.4G	/home/ubuntu/.carbon/data/application.db
    408G	/home/ubuntu/.carbon/data/blockstore.db
    334G	/home/ubuntu/.carbon/data/state.db
    1.6T	/home/ubuntu/.carbon/data/tx_index.db
    ```

    ```bash
    # after
    7.4G	/home/ubuntu/.carbon/rocksdata/application.db
    468G	/home/ubuntu/.carbon/rocksdata/blockstore.db
    1.4T	/home/ubuntu/.carbon/rocksdata/state.db
    1.9T	/home/ubuntu/.carbon/rocksdata/tx_index.db
    ```

    Create rocksdb data directory
    ```bash
    mkdir ~/.carbon/rocksdata
    ```

4. Convert leveldb to rocksdb


    ```bash

    nohup ./carbon-db-convert -dbDir="$HOME/.carbon/data/application.db" -outDir="$HOME/.carbon/rocksdata" > convert_application.out 2>&1 &

    nohup ./carbon-db-convert -dbDir="$HOME/.carbon/data/blockstore.db" -outDir="$HOME/.carbon/rocksdata" > convert_blockstore.out 2>&1 &

    nohup ./carbon-db-convert -dbDir="$HOME/.carbon/data/state.db" -outDir="$HOME/.carbon/rocksdata" > convert_state.out 2>&1 &

    nohup ./carbon-db-convert -dbDir="$HOME/.carbon/data/tx_index.db" -outDir="$HOME/.carbon/rocksdata" > convert_tx_index.out 2>&1 &
    ```

    ```bash
    # estimated progress using logs

    # convert_application.out
    Count: 7730000
    # convert_blockstore.out
    Count: 140560000
    # convert_state.out
    Count: 84060000
    # convert_tx_index.out
    Count: 9283320000
    ```

5. Ensure process has completed and no errors in logs
   ```bash
   ps aux | grep convert_
   tail convert_*.out
   ```

6. Move /rocksdbdata to /data
    ```bash
    rm -rf ~/.carbon/data/application.db && mv ~/.carbon/rocksdata/application.db ~/.carbon/data
    rm -rf ~/.carbon/data/blockstore.db && mv ~/.carbon/rocksdata/blockstore.db ~/.carbon/data
    rm -rf ~/.carbon/data/state.db && mv ~/.carbon/rocksdata/state.db ~/.carbon/data
    rm -rf ~/.carbon/data/tx_index.db && mv ~/.carbon/rocksdata/tx_index.db ~/.carbon/data
    ```
