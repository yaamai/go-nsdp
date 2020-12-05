package nsdp

import (
	"bytes"
)

type Body []TLV

func (b Body) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	err := b.MarshalBinaryBufferWithOption(&buf, false)
	return buf.Bytes(), err
}

func (b *Body) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)
	return b.UnmarshalBinaryBuffer(r)
}

func (b Body) MarshalBinaryBuffer(buf *bytes.Buffer) error {
	return b.MarshalBinaryBufferWithOption(buf, true)
}

func (b Body) MarshalBinaryBufferWithOption(buf *bytes.Buffer, skipValue bool) error {
	for _, tlv := range b {
		tag := tlv.Tag()

		length := uint16(0)
		value := []byte{}
		if !skipValue {
			length = tlv.Length()
			value = tlv.Value()
		}

		buf.WriteByte(byte(tag >> 8))
		buf.WriteByte(byte(tag & 0xff))
		buf.WriteByte(byte(length >> 8))
		buf.WriteByte(byte(length & 0xff))
		buf.Write(value)
	}

	return nil
}

func (b *Body) UnmarshalBinaryBuffer(r *bytes.Reader) error {
	for r.Len() > 4 {
		tag := readUint16(r)
		length := readUint16(r)

		if r.Len() < int(length) {
			break
		}

		value := make([]byte, length)
		r.Read(value)

		*b = append(*b, NewTLVFromBytes(tag, length, value))
	}

	return nil
}
