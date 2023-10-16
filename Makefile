build: go.sum
	go build -tags rocksdb -o db-convert ./main.go

install: ./db-convert
	cp db-convert ~/bin