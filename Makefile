SUFFIX=

GO=go
OUTDIR=.

.PHONY: all test clean huffmancmd lz77cmd

all: huffmancmd lz77cmd

huffmancmd:
	$(GO) build -o $(OUTDIR)/huffmancmd ./cmd/huffman

lz77cmd:
	$(GO) build -o $(OUTDIR)/lz77cmd ./cmd/lz77

test:
	$(GO) test \
	  -cover `go list -f '{{if ne .Name "main"}} {{.ImportPath}} {{end}}' ./...` \
	  -coverprofile cover.out

clean:
	-rm $(OUTDIR)/huffmancmd $(OUTDIR)/lz77cmd $(OUTDIR)/cover.out
