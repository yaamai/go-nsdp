package nsdp

import (
	"bytes"
)

type Marker struct {
	EndOfData [4]byte
}

func (m Marker) WriteToBuffer(b *bytes.Buffer) {
	b.Write(m.EndOfData[:])
}

func (m *Marker) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 4 {
		return
	}
	b.Read(m.EndOfData[:])
}

func (b Marker) String() string {
	if bytes.Compare(b.EndOfData[:], DefaultMarker.EndOfData[:]) == 0 {
		return "<MARK>"
	} else {
		return "<INVALID-MARK>"
	}
}
