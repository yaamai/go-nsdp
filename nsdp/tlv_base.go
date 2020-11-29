package nsdp

import (
	"bytes"
)

type NSDPTLV interface {
	String() string
	WriteToBuffer(b *bytes.Buffer)
	ReadFromBuffer(b *bytes.Reader)
}

type NSDPTLVBase struct {
	Tag    uint16
	Length uint16
	Value  []byte
}
