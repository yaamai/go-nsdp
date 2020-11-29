package nsdp

import (
	"bytes"
)

type Msg struct {
	Header
	Body
	Marker
}

func ParseMsg(buf []byte) *Msg {
	m := &Msg{}
	r := bytes.NewReader(buf)
	m.Header.ReadFromBuffer(r)
	m.Body.ReadFromBuffer(r)
	m.Marker.ReadFromBuffer(r)
	return m
}

func (m Msg) WriteToBuffer(b *bytes.Buffer) {
	m.Header.WriteToBuffer(b)
	m.Body.WriteToBuffer(b)
	m.Marker.WriteToBuffer(b)
}

func (m Msg) Bytes() []byte {
	b := bytes.Buffer{}
	m.WriteToBuffer(&b)
	return b.Bytes()
}
