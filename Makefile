.PHONY: build run migrate

build:
	go build -ldflags="-s -w" -o bin/cooper main.go

run:
	go run main.go

setup:
	@mkdir -p storage/db/
	@test -f storage/db/cooper.db || touch storage/db/cooper.db

dev: setup run

clean:
	rm -f bin/cooper
	rm -f cooper.db
