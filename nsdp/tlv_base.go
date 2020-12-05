package nsdp

type TLV interface {
	Tag() uint16
	Length() uint16
	Value() []byte
}

func ParseTLVs(tag uint16, length uint16, value []byte) TLV {
	switch tag {
	case 0x0001:
		return &ModelName{StringValue(value)}
	case 0x0003:
		return &HostName{StringValue(value)}
	case 0x0017:
		return &AuthV2PasswordSalt{BytesValue(value)}
	case 0x001a:
		return &AuthV2Password{BytesValue(value)}
	case 0x0c00:
		return (&PortLinkStatus{}).FromBytes(value)
	case 0x1000:
		return (&PortStatistics{}).FromBytes(value)
	case 0x2400:
		return (&PortVlanMembers{}).FromBytes(value)
	case 0x2800:
		return (&TagVlanMembers{}).FromBytes(value)
	}
	return nil
}
