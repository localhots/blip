.PHONY: test bench

test:
	go test -v ./log

bench:
	go test -v ./log -bench=. -benchmem -run=Benchmark
