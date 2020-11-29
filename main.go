package main

import (
	"log"
	"net"
	"time"
)

// get first non-loopback address, intf-name and mac
func getSelfIntfAndIp() (string, string, []byte, error) {
	intfs, err := net.Interfaces()
	if err != nil {
		return "", "", nil, err
	}

	for _, intf := range intfs {
		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), intf.Name, intf.HardwareAddr, nil
				}
			}
		}
	}

	return "", "", nil, nil
}

type NSDPClient struct {
	intfName   string
	intfHwAddr []byte
	conn       *net.UDPConn
}

func NewNSDPClient() (*NSDPClient, error) {
	selfAddrStr, intfName, intfHwAddr, err := getSelfIntfAndIp()
	if err != nil {
		return nil, err
	}

	selfAddr, err := net.ResolveUDPAddr("udp", selfAddrStr+":63321")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", selfAddr)
	if err != nil {
		return nil, err
	}

	return &NSDPClient{conn: conn, intfHwAddr: intfHwAddr, intfName: intfName}, nil
}

func (c *NSDPClient) send() error {
	sendToAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63322")
	if err != nil {
		return err
	}

	queryModel := []byte{
		0x01,
		0x01,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x01,
		0x4E, 0x53, 0x44, 0x50,
		0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00,
		//		0x00, 0x00, 0xff, 0xff,
	}

	for idx, b := range c.intfHwAddr {
		queryModel[8+idx] = b
	}

	go func() {
		for {
			buf := make([]byte, 65535)
			readLen, _, err := c.conn.ReadFrom(buf)
			log.Println(readLen, buf[:readLen], err)

		}
	}()

	go func() {
		recvAllAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63321")
		if err != nil {
			log.Println(err)
		}

		conn2, err := net.ListenUDP("udp", recvAllAddr)
		if err != nil {
			log.Println(err)
		}
		buf := make([]byte, 65535)
		readLen, _, err := conn2.ReadFrom(buf)
		if err != nil {
			log.Println(err)
		}
		log.Println(readLen, buf[:readLen], err)
	}()

	for idx := 0; idx < 120; idx++ {
		writeLen, err := c.conn.WriteTo(queryModel, sendToAddr)
		log.Println(writeLen, err)
		time.Sleep(3 * time.Second)
	}

	return nil
}

func main() {
	msg := NSDPMsg{}
	log.Println(msg.Result)
	c, err := NewNSDPClient()
	if err != nil {
		log.Fatalln(err)
	}
	c.send()

	for {
	}
}
