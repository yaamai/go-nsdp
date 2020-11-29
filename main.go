package main

import (
	"gs308e/nsdp"
	"log"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	tlv := &nsdp.TLVBase{Tag: 1, Length: 0}
	log.Println(c.Read([]nsdp.TLV{tlv}))

	for {
	}
}
