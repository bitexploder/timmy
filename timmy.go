package main

import (
	"fmt"
	"net"
)

type Mitmer struct {
	Conn *net.TCPConn
}

func (m *Mitmer) MitmConn() {
	origAddr, err := GetOriginalDST(m.Conn)
	if err != nil {
		fmt.Println("get orig addr err: ", err)
		return
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

		m := Mitmer{Conn: conn.(*net.TCPConn)}
		go m.MitmConn()
	}
}
