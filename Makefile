SUFFIX=

GO=go
OUTDIR=.

.PHONY: all test clean huffmancmd

all: huffmancmd

huffmancmd:
	$(GO) build -o $(OUTDIR)/huffmancmd ./cmd/huffman

test:
	$(GO) test \
	  -cover `go list -f '{{if ne .Name "main"}} {{.ImportPath}} {{end}}' ./...` \
	  -coverprofile cover.out

clean:
	-rm $(OUTDIR)/huffmancmd $(OUTDIR)/cover.out
