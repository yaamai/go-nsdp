package nsdp

type StringValue string

func (t StringValue) Length() uint16 {
	return uint16(len(t))
}

func (t StringValue) Value() []byte {
	return []byte(t)
}

func (t StringValue) String() string {
	return string(t)
}

type BytesValue []byte

func (t BytesValue) Length() uint16 {
	return uint16(len(t))
}

func (t BytesValue) Value() []byte {
	return t
}

func (t BytesValue) String() string {
	return string(t)
}
