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
	tlv := &nsdp.TLVModelName{}
	log.Println(c.Read([]nsdp.TLV{tlv}))

	for {
	}
}
