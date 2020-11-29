package nsdp

import (
	"bytes"
	"fmt"
)

type NSDPTLVUnknown struct {
	NSDPTLVBase
}

func (t NSDPTLVUnknown) WriteToBuffer(b *bytes.Buffer) {
	b.WriteByte(byte(t.Tag >> 8))
	b.WriteByte(byte(t.Tag & 0xff))
	b.WriteByte(byte(t.Length >> 8))
	b.WriteByte(byte(t.Length & 0xff))
	b.Write(t.Value)
}
func (t *NSDPTLVUnknown) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 4 {
		return
	}
	t.Tag = uint16(readInt16(b))
	t.Length = uint16(readInt16(b))

	if b.Len() < int(t.Length) {
		return
	}
	t.Value = make([]byte, t.Length)
	b.Read(t.Value)
}
func (t NSDPTLVUnknown) String() string {
	return fmt.Sprintf("T: %d, V: %v(%d)", t.Tag, t.Value, t.Length)
}
