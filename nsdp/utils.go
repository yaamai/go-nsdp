package nsdp

import (
	"bytes"
)

var (
	EmptyMac      = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	DefaultHeader = Header{Version: 1, Op: 1, HostMac: EmptyMac, DeviceMac: EmptyMac, Seq: 1, Signature: [4]byte{0x4E, 0x53, 0x44, 0x50}}
	DefaultBody   = Body{}
	DefaultMarker = Marker{EndOfData: [4]byte{0xff, 0xff, 0x00, 0x00}}
	DefaultMsg    = Msg{DefaultHeader, DefaultBody, DefaultMarker}
)

func readInt8(b *bytes.Reader) int8 {
	v, _ := b.ReadByte()
	return int8(v)
}
func readInt16(b *bytes.Reader) int16 {
	v1, _ := b.ReadByte()
	v2, _ := b.ReadByte()

	var v int16
	v = int16(v1) << 8
	v = v | int16(v2)
	return v
}

func parsePortsBits(buf []byte) []int {
	ports := []int{}
	portNum := 1
	for _, b := range buf {
		for idx := 0x80; idx > 0; idx = idx >> 1 {
			if b&byte(idx) != 0 {
				ports = append(ports, portNum)
			}
			portNum += 1
		}
	}

	return ports
}

func combinePortsBits(ports []int) []byte {
	// determine output byte length
	max := 0
	for _, p := range ports {
		if max < p {
			max = p
		}
	}
	bufLen := ((max - 1) / 8) + 1
	buf := make([]byte, bufLen)

	for _, p := range ports {
		bytePos := ((p - 1) / 8)
		bitPos := (p - 1) % 8
		buf[bytePos] |= 0x80 >> bitPos
	}
	return buf
}
