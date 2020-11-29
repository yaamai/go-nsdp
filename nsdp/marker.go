package nsdp

import (
	"bytes"
)

type NSDPMarker struct {
	EndOfData [4]byte
}

func (m NSDPMarker) WriteToBuffer(b *bytes.Buffer) {
	b.Write(m.EndOfData[:])
}

func (m *NSDPMarker) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 4 {
		return
	}
	b.Read(m.EndOfData[:])
}

func (b NSDPMarker) String() string {
	if bytes.Compare(b.EndOfData[:], NSDPDefaultMarker.EndOfData[:]) == 0 {
		return "<MARK>"
	} else {
		return "<INVALID-MARK>"
	}
}
