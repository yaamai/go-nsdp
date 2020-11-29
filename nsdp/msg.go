package nsdp

import (
	"bytes"
)

type NSDPMsg struct {
	NSDPHeader
	NSDPBody
	NSDPMarker
}

func ParseNSDPMsg(buf []byte) *NSDPMsg {
	m := &NSDPMsg{}
	r := bytes.NewReader(buf)
	m.NSDPHeader.ReadFromBuffer(r)
	m.NSDPBody.ReadFromBuffer(r)
	m.NSDPMarker.ReadFromBuffer(r)
	return m
}

func (m NSDPMsg) WriteToBuffer(b *bytes.Buffer) {
	m.NSDPHeader.WriteToBuffer(b)
	m.NSDPBody.WriteToBuffer(b)
	m.NSDPMarker.WriteToBuffer(b)
}

func (m NSDPMsg) Bytes() []byte {
	b := bytes.Buffer{}
	m.WriteToBuffer(&b)
	return b.Bytes()
}
