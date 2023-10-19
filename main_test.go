package main

import (
	"fmt"
	"testing"
	_ "github.com/joho/godotenv/autoload" // reads .env file
)

// set the .carbon/data paths of all the db files to use for the iterator benchmark (replace with your own dirs!)
var pathTable = []struct {
	lvlDbPath string
	dbName string
} {
	{lvlDbPath: ".carbon/data/application.db", dbName: "application"},
	{lvlDbPath: ".carbon/data/blockstore.db", dbName: "blockstore"},
	{lvlDbPath: ".carbon/data/evidence.db", dbName: "evidence"},
	{lvlDbPath: ".carbon/data/snapshots/metadata.db", dbName: "metadata"},
	{lvlDbPath: ".carbon/data/state.db", dbName: "state"},
	{lvlDbPath: ".carbon/data/tx_index.db", dbName: "tx_index"},
}

// Run: go test -bench=. -benchmem -tags rocksdb
func BenchmarkIterateDb(b *testing.B) {
	// set the output directory for the generated db files (replace with your own dirs!)
	const outputDir = "db-convert/testing/output/"
	for _, v := range pathTable {
		b.Run(fmt.Sprintf("db-%s", v.dbName), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                iterateDb(v.lvlDbPath, v.dbName, outputDir)
            }
        })
	}
}