// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

package main

import (
    "encoding/binary"
    "net"
)

type Key struct {
    SrcAddr string `validate:"required"`
    DstPort uint16`validate:"required"`
    Proto   uint8 `validate:"required"`
}

func (k *Key) GetBytes() []byte {
    res := make([]byte, 0, 8)

    dstPortBytes := make([]byte, 2)
    binary.BigEndian.PutUint16(dstPortBytes, k.DstPort)

    res = append(res, net.ParseIP(k.SrcAddr).To4()...)
    res = append(res, dstPortBytes...)
    res = append(res, k.Proto)
    res = append(res, byte(0))

    return res
}
