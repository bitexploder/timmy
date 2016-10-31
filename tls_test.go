package main

import (
	"crypto/tls"
	"fmt"
	"testing"
)

var tc *tls.Config = &tls.Config{InsecureSkipVerify: true}

func TestTls(t *testing.T) {
	conn, err := tls.Dial("tcp", "127.0.0.1:9090", tc)
	if err != nil {
		fmt.Println("Error connecting...")
	}
	conn.Write([]byte("Hello, TLS\n"))
	conn.Close()
}
