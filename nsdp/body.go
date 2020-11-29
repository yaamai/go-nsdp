package nsdp

import (
	"bytes"
)

type Body struct {
	Body []TLV
}

func (b Body) WriteToBuffer(buf *bytes.Buffer) {
	for idx, _ := range b.Body {
		b.Body[idx].WriteToBuffer(buf)
	}
}
func (b *Body) ReadFromBuffer(buf *bytes.Reader) {
	for buf.Len() > 4 {
		tlv := TLVBase{}
		tlv.ReadFromBuffer(buf)
		b.Body = append(b.Body, &tlv)
	}
}
