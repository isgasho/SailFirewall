// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

package main

import (
    "flag"
    "fmt"
    "github.com/dropbox/goebpf"
    "go.uber.org/zap"
    "os"
    "os/signal"
)

var (
    logger *zap.Logger
    aclist goebpf.Map
    iface = flag.String("iface", "", "Interface to bind XDP program to")
    elf = flag.String("elf", "ebpf_prog/xdp_fw.elf", "clang/llvm compiled binary file")
    apiAddr = flag.String("addr", ":8765", "api server bind address")
)

func init() {
    logger, _ = zap.NewProduction()
}

func main() {
	flag.Parse()
	if *iface == "" {
		logger.Fatal("-iface is required")
	}

	// Create eBPF system
	bpf := goebpf.NewDefaultEbpfSystem()
	// Load .ELF files compiled by clang/llvm
	err := bpf.LoadElf(*elf)
	if err != nil {
		logger.Fatal("LoadElf() failed", zap.Error(err))
	}
	printBpfInfo(bpf)

	// Get eBPF maps
	aclist = bpf.GetMapByName("aclist")
	if aclist == nil {
		logger.Fatal("eBPF map 'aclist' not found")
	}

	// Get XDP program. Name simply matches function from xdp_fw.c:
	//      int firewall(struct xdp_md *ctx) {
	xdp := bpf.GetProgramByName("firewall")
	if xdp == nil {
		logger.Fatal("Program 'firewall' not found")
	}

    key := &Key{
        SrcAddr: "127.0.0.1",
        Proto: 0x06,
        DstPort: 8765,
    }

    err = aclist.Insert(key.GetBytes(), 0)
    if err != nil {
        logger.Fatal("add rule for api server failed", zap.Error(err))
    }

	fmt.Println()

	// Load XDP program into kernel
	err = xdp.Load()
	if err != nil {
		logger.Fatal("xdp.Load() error", zap.Error(err))
	}

	// Attach to interface
	err = xdp.Attach(*iface)
	if err != nil {
		logger.Fatal("xdp.Attach() error", zap.Error(err))
	}
	defer xdp.Detach()

	// Add CTRL+C handler
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

    go runApiServer()

    <-ctrlC
    logger.Info("Detaching program and exit")
}
