package nsdp

type Msg struct {
	Header
	Body
	Marker
}

func (m *Msg) UnmarshalBinary(buf []byte) error {
	m.Header.UnmarshalBinary(buf[:HeaderLength])
	m.Body.UnmarshalBinary(buf[HeaderLength : len(buf)-MarkerLength])
	m.Marker.UnmarshalBinary(buf[len(buf)-MarkerLength:])
	return nil
}

func (m Msg) MarshalBinary() ([]byte, error) {
	// if read op, write only tag (length, value not needed)
	skipValue := false
	if m.Header.Op == 1 {
		skipValue = true
	}

	header, err := m.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	body, err := m.Body.MarshalBinaryWithOption(skipValue)
	if err != nil {
		return nil, err
	}

	marker, err := m.Marker.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(header, append(body, marker...)...), nil
}
