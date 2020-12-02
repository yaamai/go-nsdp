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

	resp := c.Read(nsdp.TagVlanMembers{})
	log.Printf("%v", resp)
}
