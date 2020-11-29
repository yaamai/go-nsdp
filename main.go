package main

import (
	"log"
)

func main() {
	c, err := NewNSDPClient()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(c.Read([]NSDPTLV{NSDPTLV{Tag: 1, Length: 0}}))

	for {
	}
}
