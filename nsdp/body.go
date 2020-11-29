package nsdp

import (
	"bytes"
)

type NSDPBody struct {
	Body []NSDPTLV
}

func (b NSDPBody) WriteToBuffer(buf *bytes.Buffer) {
	for idx, _ := range b.Body {
		b.Body[idx].WriteToBuffer(buf)
	}
}
func (b *NSDPBody) ReadFromBuffer(buf *bytes.Reader) {
	for buf.Len() > 4 {
		tlv := NSDPTLVUnknown{}
		tlv.ReadFromBuffer(buf)
		b.Body = append(b.Body, &tlv)
	}
}
