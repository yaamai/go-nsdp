package main

import (
	"bytes"
	"gs308e/nsdp"
	"log"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	var cur, prev *nsdp.Msg
	for idx := 1; idx < 0xffff; idx++ {
		if idx >= 0x1000 && idx <= 0x11ff {
			continue
		}
		cur = c.Read(&nsdp.TLVBase{Tag: uint16(idx), Length: 0})
		log.Printf("%x, %v", idx, cur)

		if cur != nil && prev != nil {
			prev.Seq = cur.Seq
			prev.Unknown1 = cur.Unknown1
			// log.Println(cur.Bytes())
			// log.Println(prev.Bytes())
			if bytes.Compare(cur.Bytes(), prev.Bytes()) == 0 {
				log.Printf("skip")
				idx += 0x80 - 1
				prev = nil
				continue
			}
		}
		prev = cur
	}

	for {
	}
}
