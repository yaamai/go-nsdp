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
