.PHONY: test bench

test:
	go test -v .

bench:
	go test -v . -bench=. -benchmem -run=Benchmark
