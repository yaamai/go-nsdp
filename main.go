package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gs308e/nsdp"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	data := map[string]string{}

	dataBytes, err := ioutil.ReadFile("dump.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(dataBytes, &data)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		jsonBytes, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Fatalln(err)
		}

		ioutil.WriteFile("dump.json", jsonBytes, 0644)
		os.Exit(0)
	}()

	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	var cur, prev *nsdp.Msg
	for idx := 1; idx < 0xffff; idx++ {
		if _, ok := data[fmt.Sprintf("%x", idx)]; ok {
			continue
		}
		if idx >= 0x1000 && idx <= 0x13ff {
			continue
		}
		cur = c.Read(&nsdp.TLVBase{Tag: uint16(idx), Length: 0})
		log.Printf("%x, %v", idx, cur)
		data[fmt.Sprintf("%x", idx)] = fmt.Sprintf("%v", cur)

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
