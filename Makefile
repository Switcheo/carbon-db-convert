build: go.sum
	go build -tags rocksdb -o carbon-db-convert ./main.go

install: ./carbon-db-convert
	cp carbon-db-convert ~/bin