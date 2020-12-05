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
	buf := bytes.Buffer{}
	err := m.MarshalBinaryBuffer(&buf)
	return buf.Bytes(), err
}
func (m *Marker) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)
	return m.UnmarshalBinaryBuffer(r)
}

func (m Marker) MarshalBinaryBuffer(buf *bytes.Buffer) error {
	buf.Write(m.EndOfData[:])
	return nil
}

func (m *Marker) UnmarshalBinaryBuffer(r *bytes.Reader) error {
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
