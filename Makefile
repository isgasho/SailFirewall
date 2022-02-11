# Copyright (c) 2022 Hevienz.
# Full license can be found in the LICENSE file.

GOCMD := go
GOBUILD := $(GOCMD) build -ldflags '-s -w'
GOCLEAN := $(GOCMD) clean

CLANG := clang
CLANG_INCLUDE := -I../../dropbox/goebpf

GO_BINARY := SailFirewall

EBPF_SOURCE := ebpf_prog/xdp_fw.c
EBPF_BINARY := ebpf_prog/xdp_fw.elf

all: build_bpf build_go

build_bpf: $(EBPF_BINARY)

build_go:
	$(GOBUILD) -v

clean:
	$(GOCLEAN)
	rm -f $(GO_BINARY)
	rm -f $(EBPF_BINARY)

$(EBPF_BINARY): $(EBPF_SOURCE)
	$(CLANG) $(CLANG_INCLUDE) -O2 -target bpf -c $^  -o $@
