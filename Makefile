proj:= costool
bindir := /usr/local/bin
.PHONY: build
build:
	go build -o bin/$(proj) dacazh.com/$(proj)

.PHONY: run
run:
	@go run dacazh.com/$(proj)

.PHONY: clean test install uninstall
clean:
	rm bin/*
test:
	@make
	# rm test -rf
	# ./bin/$(proj)  -p README.md bin/
	# ./bin/$(proj) -g eset test/

install:
	install -d $(bindir)
	install -m 0755 bin/$(proj)  $(bindir)

uninstall:
	rm  $(bindir)/$(proj)