
package main

import (
	"fmt"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/joho/godotenv/autoload" // reads .env file

	cdb "github.com/cometbft/cometbft-db"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	o := &opt.Options{
		// The default value is nil
		Filter: filter.NewBloomFilter(10),
		// Use 1 GiB instead of default 8 MiB
		BlockCacheCapacity: opt.GiB,
		// Use 64 MiB instead of default 4 MiB
		WriteBuffer:                           64 * opt.MiB,
		CompactionTableSize:                   8 * opt.MiB,
		CompactionTotalSize:                   40 * opt.MiB,
		CompactionTotalSizeMultiplierPerLevel: []float64{1, 1, 10, 100, 1000, 10000, 100000},
		
		DisableSeeksCompaction: true,
	}

	// get machine home directory
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// default db directory: $HOME/.carbon/data/
	dbDir := filepath.Join(home, ".carbon", "data")
	// get root directory of this project
	_, b, _, _ := runtime.Caller(0)
	// default output directory: $ROOT/output/
	outputDir := filepath.Join(filepath.Dir(b), "/output")

	dbDirFlag := flag.String("dbDir", dbDir, "directory containing leveldb files")
	outputDirFlag := flag.String("outDir", outputDir, "directory for generated db files")

	if len(os.Args[1:]) < 2 {
		fmt.Printf("Using default values for the db and/or output directories.\n")
		fmt.Printf("Specify using : db-convert -dbDir= <directory> -outDir= <directory>\n\n")
	}

	// parse flag string pointers
	flag.Parse()
	dbDir = *dbDirFlag
	outputDir = *outputDirFlag

	createErr := os.MkdirAll(outputDir, os.ModePerm)
	if createErr != nil {
		panic(createErr)
	}
	fmt.Printf("output dir: %s\n", outputDir)

	// walk through and open leveldb files in the .carbon/data directory
	fmt.Printf("db dir: %s\n\n", dbDir)
	fileErr := filepath.WalkDir(dbDir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		// opens the leveldb file
		if filepath.Ext(path) == ".db" {
			fmt.Printf("path: %s | db: %s\n", path, dir.Name())

			dbName := strings.TrimSuffix(dir.Name(), ".db")
			lvlDb, dbErr := cdb.NewGoLevelDBWithOpts(dbName, filepath.Dir(path), o)

			if dbErr != nil {
				panic(dbErr)
			}

			// create a rocksdb file with the same name in the output directory
			fmt.Println("Creating rocksdb file...")
			rocksDb, newDbErr := cdb.NewDB(dbName, cdb.RocksDBBackend, outputDir)

			if newDbErr != nil {
				panic(newDbErr)
			}

			itr, itrErr := lvlDb.Iterator(nil, nil)

			if itrErr != nil {
				panic(itrErr)
			}

			count := 0
			for ; itr.Valid(); itr.Next() {
				key := itr.Key()
				value := itr.Value()

				err := rocksDb.Set(key, value)

				if err != nil {
					panic(err)
				}

				count++

				if count%10000 == 0 {
					fmt.Printf("Count: %d\n", count)
					runtime.GC() // Force GC
				}
			}

			itr.Close()
			rocksDb.Close()
			lvlDb.Close()
		}
		return nil
	})
	if fileErr != nil {
		println("filepath WalkDir error!")
		panic(fileErr)
	}
}

// version of the main iteration for use in testing/benching
func iterateDb(lvlDbPath string, dbName string, outputDir string) {
	
	o := &opt.Options{
		// The default value is nil
		Filter: filter.NewBloomFilter(10),
		// Use 1 GiB instead of default 8 MiB
		BlockCacheCapacity: opt.GiB,
		// Use 64 MiB instead of default 4 MiB
		WriteBuffer:                           64 * opt.MiB,
		CompactionTableSize:                   8 * opt.MiB,
		CompactionTotalSize:                   40 * opt.MiB,
		CompactionTotalSizeMultiplierPerLevel: []float64{1, 1, 10, 100, 1000, 10000, 100000},
		
		DisableSeeksCompaction: true,
	}

	lvlDb, dbErr := cdb.NewGoLevelDBWithOpts(dbName, filepath.Dir(lvlDbPath), o)

	if dbErr != nil {
		panic(dbErr)
	}

	// create a rocksdb file with the same name in the output directory
	rocksDb, newDbErr := cdb.NewDB(dbName, cdb.RocksDBBackend, outputDir)

	if newDbErr != nil {
		panic(newDbErr)
	}

	itr, itrErr := lvlDb.Iterator(nil, nil)

	if itrErr != nil {
		panic(itrErr)
	}

	count := 0
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		value := itr.Value()

		err := rocksDb.Set(key, value)

		if err != nil {
			panic(err)
		}

		count++

		if count%10000 == 0 {
			// fmt.Printf("Count: %d\n", count)
			runtime.GC() // Force GC
		}
	}

	itr.Close()
	rocksDb.Close()
	lvlDb.Close()
}