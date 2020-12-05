package nsdp

import (
	"bytes"
	"errors"
)

const (
	MarkerLength = 4
)

type Marker struct {
	EndOfData [4]byte
}

func (m Marker) MarshalBinary() ([]byte, error) {
	return m.EndOfData[:], nil
}

func (m *Marker) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)
	if r.Len() < 4 {
		return errors.New("too short end of tlv marker")
	}
	r.Read(m.EndOfData[:])

	return nil
}

func (b Marker) String() string {
	if bytes.Compare(b.EndOfData[:], DefaultMarker.EndOfData[:]) == 0 {
		return "<END>"
	} else {
		return "<INVALID-END>"
	}
}
