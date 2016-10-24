package main

import (
	"fmt"
	"net"
)

type Mitmer struct {
	InConn  net.Conn
	OutConn net.Conn
}

func (m *Mitmer) MitmConn() {
	origAddr, err := GetOriginalDST(m.InConn.(*net.TCPConn))
	if err != nil {
		fmt.Println("get orig addr err: ", err)
		return
	}

	outc, err := net.Dial("tcp", fmt.Sprintf("%+v", origAddr))
	m.OutConn = outc

	if err != nil {
		fmt.Println("err connecting to orig dst: ", err)
		return
	}

	// Clean up
	defer m.InConn.Close()
	defer outc.Close()

	go func() {
		for {
			b := make([]byte, 1024)
			n, err := outc.Read(b)
			if err != nil {
				fmt.Println("err reading victim dest: ", err)
				break
			}
			n, err = m.InConn.Write(b[:n])
			//fmt.Printf("Writing back to victim: %+v\n", b[:n])
			if err != nil {
				fmt.Println("err writing back to victim: ", err)
				break
			}

		}
	}()

	// Set up the victim->server data pump
	for {
		b := make([]byte, 1024)
		n, err := m.InConn.Read(b)
		//fmt.Printf("Read bytes[%d] from [%+v]\n", n, origAddr)

		if err != nil {
			fmt.Println("err reading victim: ", err)
			break
		}
		n, err = outc.Write(b[:n])
		//fmt.Printf("Writing: %+v\n", b[:n])
		if err != nil {
			fmt.Println("err writing victim dest: ", err)
			break
		}

	}

	fmt.Printf("orig: %+v\n", origAddr)
}

func main() {
	fmt.Println("Timmy starting up")
	//	s, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 20755})
	s, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 20755})
	if err != nil {
		fmt.Println("listen err: ", err)
		return
	}

	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			return
		}

		m := Mitmer{InConn: conn.(*net.TCPConn)}
		go m.MitmConn()
	}
}
