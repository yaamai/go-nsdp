package main

import (
	"gs308e/nsdp"
	"log"
	"math/rand"
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
	anyAddr    *net.UDPAddr
	intfName   string
	intfHwAddr []byte
	conn       *net.UDPConn
	seq        uint16
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

	anyAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63322")
	if err != nil {
		return nil, err
	}

	// to avoid ignore msg, set random sequence number
	rand.Seed(time.Now().UnixNano())
	seq := uint16(rand.Intn(0xffff))
	log.Println(seq)

	return &NSDPClient{anyAddr: anyAddr, conn: conn, intfHwAddr: intfHwAddr, intfName: intfName, seq: seq}, nil
}

func (c *NSDPClient) SendRecvMsg(msg nsdp.NSDPMsg) *nsdp.NSDPMsg {
	recvCh := make(chan bool, 1)
	buf := make([]byte, 65535)
	readLen := 0
	go func() {
		readLen, _, _ = c.conn.ReadFrom(buf)
		recvCh <- true
	}()

	retry := 0
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for retry < 3 {
		select {
		case <-recvCh:
			log.Println(readLen, buf[:readLen])
			return nsdp.ParseNSDPMsg(buf[:readLen])
		case <-ticker.C:
			writeLen, err := c.conn.WriteTo(msg.Bytes(), c.anyAddr)
			log.Println(writeLen, err)
			retry += 1
		}
	}

	return nil
}

func (c *NSDPClient) Read(msg []nsdp.NSDPTLV) *nsdp.NSDPMsg {
	m := nsdp.NSDPMsg(nsdp.NSDPDefaultMsg)
	m.Op = 1
	m.Seq = c.seq
	m.HostMac = c.intfHwAddr
	m.Body = msg

	return c.SendRecvMsg(m)
}

func (c *NSDPClient) Write(msg []nsdp.NSDPTLV) *nsdp.NSDPMsg {
	m := nsdp.NSDPMsg(nsdp.NSDPDefaultMsg)
	m.Op = 3
	m.Seq = c.seq
	m.HostMac = c.intfHwAddr
	m.Body = msg

	return c.SendRecvMsg(m)
}
