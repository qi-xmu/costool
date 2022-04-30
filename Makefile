proj:= costool

.PHONY: build
build:
	go build -trimpath -o bin/$(proj)

.PHONY: run
run:
	@go run dacazh.com/$(proj)

.PHONY: clean test
clean:
	rm bin/*
test:
	@go run dacazh.com/$(proj) -l .