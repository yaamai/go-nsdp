package nsdp

import (
	"bytes"
)

type Msg struct {
	Header
	Body
	Marker
}

func (m Msg) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	err := m.MarshalBinaryBuffer(&buf)
	return buf.Bytes(), err
}

func (m *Msg) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)
	return m.UnmarshalBinaryBuffer(r)
}

func (m *Msg) UnmarshalBinaryBuffer(buf *bytes.Reader) error {
	err := m.Header.UnmarshalBinaryBuffer(buf)
	if err != nil {
		return err
	}
	err = m.Body.UnmarshalBinaryBuffer(buf)
	if err != nil {
		return err
	}
	err = m.Marker.UnmarshalBinaryBuffer(buf)
	if err != nil {
		return err
	}
	return nil
}

func (m Msg) MarshalBinaryBuffer(buf *bytes.Buffer) error {
	// if read op, write only tag (length, value not needed)
	skipValue := false
	if m.Header.Op == 1 {
		skipValue = true
	}

	err := m.Header.MarshalBinaryBuffer(buf)
	if err != nil {
		return err
	}
	err = m.Body.MarshalBinaryBufferWithOption(buf, skipValue)
	if err != nil {
		return err
	}

	err = m.Marker.MarshalBinaryBuffer(buf)
	if err != nil {
		return err
	}

	return nil
}
