proj:= costool
bindir := /usr/local/bin
.PHONY: build
build:
	go build -o bin/ dacazh.com/costool 

.PHONY: run
run:
	@go run dacazh.com/$(proj)

.PHONY: clean test
clean:
	rm bin/*
test:
	@make
	# rm test -rf
	# ./bin/$(proj)  -p README.md bin/
	./bin/$(proj) -g eset test/

install:
	install -d $(bindir)
	install -m 0755 bin/$(proj)  $(bindir)