package nsdp

import (
	"bytes"
	"fmt"
)

type TLV interface {
	String() string
	WriteToBuffer(b *bytes.Buffer)
	ReadFromBuffer(b *bytes.Reader)
}

type TLVBase struct {
	Tag    uint16
	Length uint16
	Value  []byte
}

func (t TLVBase) WriteToBuffer(b *bytes.Buffer) {
	b.WriteByte(byte(t.Tag >> 8))
	b.WriteByte(byte(t.Tag & 0xff))
	b.WriteByte(byte(t.Length >> 8))
	b.WriteByte(byte(t.Length & 0xff))
	b.Write(t.Value)
}
func (t *TLVBase) ReadFromBuffer(b *bytes.Reader) {
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
func (t TLVBase) String() string {
	return fmt.Sprintf("T: %d, V: %v(%d)", t.Tag, t.Value, t.Length)
}
