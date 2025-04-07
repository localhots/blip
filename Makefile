.PHONY: test bench pprof

test:
	go test -v .

bench:
	go test -v . -bench=. -benchmem -run=Benchmark

fuzz:
	go test -v -fuzz=Fuzz -run=Fuzz -fuzztime=10s

pprof:
	go build -o /tmp/blip prof/main.go
	/tmp/blip -cpuprofile=/tmp/blip.prof
	go tool pprof /tmp/blip /tmp/blip.prof
