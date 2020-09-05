SUFFIX=

GO=go
GOLINT=golint
OUTDIR=.

ITERATIONS=5

.PHONY: all test clean huffmancmd lz77cmd lint perftestrunner perf-report gendocs

all: huffmancmd lz77cmd

huffmancmd:
	$(GO) build -o $(OUTDIR)/huffmancmd ./cmd/huffman

lz77cmd:
	$(GO) build -o $(OUTDIR)/lz77cmd ./cmd/lz77

perftestrunner:
	$(GO) build -o ./test/runner ./tools/perftestrunner

perf-report: huffmancmd lz77cmd perftestrunner
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
	@./test/runner \
	  -iters $(ITERATIONS) \
	  -cmd $(OUTDIR)/huffmancmd \
	  -workdir ./test/tmp \
	  -dir ./test/files/complexity-analysis \
	  > huffman-complexity-stats.csv
	@./test/runner \
	  -iters $(ITERATIONS) \
	  -cmd $(OUTDIR)/lz77cmd \
	  -workdir ./test/tmp \
	  -dir ./test/files/complexity-analysis \
	  > lz77-complexity-stats.csv

gendocs:
	$(GO) run ./tools/gendocs ./docs .

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
