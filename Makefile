SUFFIX=

GO=go
GOLINT=golint
OUTDIR=.

ITERATIONS=5

.PHONY: all test clean huffmancmd lz77cmd lint testrunner perf-report

all: huffmancmd lz77cmd

huffmancmd:
	$(GO) build -o $(OUTDIR)/huffmancmd ./cmd/huffman

lz77cmd:
	$(GO) build -o $(OUTDIR)/lz77cmd ./cmd/lz77

testrunner:
	$(GO) build -o ./test/runner ./test/cmd

perf-report: huffmancmd lz77cmd testrunner
	@./test/runner \
	  -iters $(ITERATIONS) \
	  -cmd $(OUTDIR)/huffmancmd \
	  -workdir ./test/tmp \
	  -dir ./test/files \
	  > huffman-stats.csv
	@./test/runner \
	  -iters $(ITERATIONS) \
	  -cmd $(OUTDIR)/lz77cmd \
	  -workdir ./test/tmp \
	  -dir ./test/files \
	  > lz77-stats.csv

test:
	@$(GO) test \
	  -coverprofile=cover.out \
	  -coverpkg=./... \
	  `go list -f '{{if ne .Name "main"}}{{.ImportPath}}{{end}}' ./...`
	@go tool cover -func=cover.out | tail -1

lint:
	@$(GO) vet ./...
	@$(GOLINT) ./...

clean:
	-rm -r \
	  $(OUTDIR)/huffmancmd \
	  $(OUTDIR)/lz77cmd \
	  ./test/runner \
	  ./test/tmp
