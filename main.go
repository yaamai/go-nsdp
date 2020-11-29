package main

import (
	"gs308e/nsdp"
	"log"
)

func main() {
	c, err := NewNSDPClient()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(c.Read([]nsdp.NSDPTLV{&nsdp.NSDPTLVUnknown{nsdp.NSDPTLVBase{Tag: 1, Length: 0}}}))

	for {
	}
}
