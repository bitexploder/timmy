package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Mitmer struct {
	Conf    Config
	InConn  net.Conn
	OutConn net.Conn
}

func (m *Mitmer) GetDest() *net.TCPAddr {
	localAddr := m.InConn.LocalAddr()
	fmt.Printf("Local Address: %+v\n", m.InConn.LocalAddr())
	s := strings.Split(localAddr.String(), ":")
	p, err := strconv.Atoi(s[1])
	if err != nil {
		fmt.Printf("GetDest: error converting port to integer")
		// return invalid address / error here?
	}
	addr := net.TCPAddr{
		net.ParseIP(s[0]),
		p,
		"",
	}

	fmt.Printf("New Addr: %+v\n", addr)

	return &addr
}
func (m *Mitmer) MitmConn() {
	m.GetDest()

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

	// Setup the server->victim data pump
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

func listener(l *net.TCPListener, c chan net.Conn) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept err: ", err)
			return
		}
		c <- conn
	}
}

func connMitmer(c chan net.Conn) {
	for {
		conn := <-c
		m := Mitmer{InConn: conn.(*net.TCPConn)}
		go m.MitmConn()
	}
}

func main() {
	fmt.Println("Timmy starting up")

	conf, err := parseFlags()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("Config: %+v\n", conf)

	listeners := make([]*net.TCPListener, 0)

	inC := make(chan net.Conn)

	for port := range conf.Ports {
		fmt.Println(port)
		l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
		if err != nil {
			fmt.Println("listen err: ", err)
			return
		}
		listeners = append(listeners, l)
	}

	for _, l := range listeners {
		go listener(l, inC)
	}
	go connMitmer(inC)
	//go listener(l, inC)
	// s, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 20755})
	// s2, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 8088})
	// if err != nil {
	// 	fmt.Println("listen err: ", err)
	// 	return
	// }

	// inC := make(chan net.Conn)
	// go listener(s, inC)
	// go listener(s2, inC)

	done := make(chan bool)
	<-done

	// for {
	// 	conn, err := s.Accept()
	// 	if err != nil {
	// 		fmt.Println("accept err:", err)
	// 		return
	// 	}

	// 	m := Mitmer{InConn: conn.(*net.TCPConn)}
	// 	go m.MitmConn()
	// }
}
