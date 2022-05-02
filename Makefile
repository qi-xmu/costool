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
	./bin/$(proj) -g hello-world
	@cat hello-world
	@rm hello-world

install:
	install -d $(bindir)
	install -m 0755 bin/$(proj)  $(bindir)

uninstall:
	rm  $(bindir)/$(proj)