.PHONY: lint test bench fuzz pprof demo demo-console demo-json

lint:
	golangci-lint run

test:
	go test -v -run=Test

bench:
	go test -v -bench=. -benchmem -run=Benchmark

fuzz:
	go test -v -fuzz=. -fuzztime=10s -run=Fuzz

pprof:
	go build -o /tmp/blip cmd/pprof/main.go
	/tmp/blip -cpuprofile=/tmp/blip.prof
	go tool pprof -http=127.0.0.1:6060 /tmp/blip /tmp/blip.prof

demo: demo-console demo-json

demo-console:
	go run cmd/demo/main.go

demo-json:
	go run cmd/demo/main.go -enc json
