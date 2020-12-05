package nsdp

type ModelName struct {
	StringValue
}

func (t ModelName) Tag() uint16 {
	return 0x0001
}

type HostName struct {
	StringValue
}

func (t HostName) Tag() uint16 {
	return 0x0003
}
