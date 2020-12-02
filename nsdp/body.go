package nsdp

import (
	"bytes"
)

type Body struct {
	Body []TLV
}

func (b Body) WriteToBuffer(buf *bytes.Buffer) {
	for _, tlv := range b.Body {
		tag := tlv.Tag()
		length := tlv.Length()
		value := tlv.Value()

		buf.WriteByte(byte(tag >> 8))
		buf.WriteByte(byte(tag & 0xff))
		buf.WriteByte(byte(length >> 8))
		buf.WriteByte(byte(length & 0xff))
		buf.Write(value)
	}
}
func (b *Body) ReadFromBuffer(buf *bytes.Reader) {
	for buf.Len() > 4 {
		tag := uint16(readInt16(buf))
		length := uint16(readInt16(buf))

		if buf.Len() < int(length) {
			break
		}

		value := make([]byte, length)
		buf.Read(value)

		b.Body = append(b.Body, ParseTLVs(tag, length, value))
	}
}
