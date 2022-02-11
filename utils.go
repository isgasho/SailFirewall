// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

package main

import (
    "fmt"
    "github.com/dropbox/goebpf"
)

// copy from https://github.com/dropbox/goebpf/blob/master/examples/xdp/basic_firewall/main.go#L119
func printBpfInfo(bpf goebpf.System) {
    fmt.Println("Maps:")
    for _, item := range bpf.GetMaps() {
        fmt.Printf("\t%s: %v, Fd %v\n", item.GetName(), item.GetType(), item.GetFd())
    }
    fmt.Println("\nPrograms:")
    for _, prog := range bpf.GetPrograms() {
        fmt.Printf("\t%s: %v, size %d, license \"%s\"\n",
            prog.GetName(), prog.GetType(), prog.GetSize(), prog.GetLicense(),
        )

    }
    fmt.Println()
}
