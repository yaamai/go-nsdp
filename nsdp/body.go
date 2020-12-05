package nsdp

import (
	"bytes"
)

type Body []TLV

func (b Body) MarshalBinary() ([]byte, error) {
	return b.MarshalBinaryWithOption(false)
}

func (b Body) MarshalBinaryWithOption(skipValue bool) ([]byte, error) {
	buf := bytes.Buffer{}

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

	return buf.Bytes(), nil
}

func (b *Body) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)

	for r.Len() > 4 {
		tag := readUint16(r)
		length := readUint16(r)

		if r.Len() < int(length) {
			break
		}

		value := make([]byte, length)
		r.Read(value)

		*b = append(*b, ParseTLVs(tag, length, value))
	}

	return nil
}
