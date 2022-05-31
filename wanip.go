package wanip

import (
	"bytes"
	"net"
	"time"
)

// Get Get wan ip
func Get() string {
	chanResult := make(chan string, 2)
	funcGet := func(dnsIP string) {
		if ip := getViaDNS(dnsIP); len(ip) > 0 {
			chanResult <- ip
		}
	}

	go funcGet("208.67.220.220")
	go funcGet("208.67.222.222")

	select {
	case ip := <-chanResult:
		return ip

	case <-time.After(time.Second * 10):
		return ""
	}
}

func getViaDNS(dnsIP string) string {
	udp, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(dnsIP),
		Port: 53,
	})

	if err != nil {
		return ""
	}
	defer udp.Close()

	// request
	if _, err := udp.Write([]byte{
		0x01, 0x12, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x6d, 0x79, 0x69,
		0x70, 0x07, 0x6f, 0x70, 0x65, 0x6e, 0x64, 0x6e, 0x73, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01,
		0x00, 0x01,
	}); err != nil {
		return ""
	}

	// receive
	buf := make([]byte, 1024)
	udp.SetDeadline(time.Now().Add(time.Second * 5))
	size, err := udp.Read(buf)

	if err != nil || size < 12 || buf[0] != 0x1 || buf[1] != 0x12 || buf[6] != 0x0 || buf[7] != 0x1 {
		return ""
	}

	// get ip
	mark := []byte{0xc0, 0x0c, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}
	pos := bytes.LastIndex(buf, mark)

	if pos != size-4-len(mark) {
		return ""
	}
	pos += len(mark)

	return (net.IP(buf[size-4 : size][:])).String()
}
