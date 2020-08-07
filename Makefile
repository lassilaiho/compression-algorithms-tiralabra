SUFFIX=

GO=go
GOLINT=golint
OUTDIR=.

.PHONY: all test clean huffmancmd lz77cmd lint

all: huffmancmd lz77cmd

huffmancmd:
	$(GO) build -o $(OUTDIR)/huffmancmd ./cmd/huffman

lz77cmd:
	$(GO) build -o $(OUTDIR)/lz77cmd ./cmd/lz77

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
	-rm $(OUTDIR)/huffmancmd $(OUTDIR)/lz77cmd $(OUTDIR)/cover.out
