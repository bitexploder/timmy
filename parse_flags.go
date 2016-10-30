package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type intslice []int

type Config struct {
	Ports    map[int]string
	TLSPorts intslice
}

// Quality of life search
func (is intslice) has(s int) bool {
	for _, i := range is {
		if i == s {
			return true
		}
	}
	return false
}

// Redir accumulator (permits multiple redir arguments)
type redir []string

func (r *redir) String() string {
	return fmt.Sprint(*r)
}

func (r *redir) Set(value string) error {
	fmt.Printf("Redir len: %d\n", len(*r))
	*r = append(*r, value)
	return nil
}

func parseFlags() (Config, error) {
	//var redir = flag.String("redir", "", "Redirect a port to a destination [orig_port]:[dest_server]:[dest_port]")
	var redirFlags redir
	var tlsPorts string
	var e error

	flag.Var(&redirFlags, "redir", "redirect port to target server:port port:serverip:serverport")
	flag.StringVar(&tlsPorts, "tlsports", "", "comma delimited list of ports to treat as TLS, ex: --tlsports 443,8088")
	flag.Parse()

	c := Config{
		Ports: make(map[int]string),
	}

	// REDIR Parsing
	for i := 0; i < len(redirFlags); i++ {
		parts := strings.Split(redirFlags[i], ":")
		if len(parts) != 3 {
			e = fmt.Errorf("--redir must have 3 arguments")
			return c, e
		}
		port, err := strconv.Atoi(parts[0])
		if err != nil {
			e = fmt.Errorf("Original port could not be converted to an integer")
			return c, e
		}

		v := net.ParseIP(parts[1])
		if v == nil {
			e = fmt.Errorf("Destination server not valid IP address")
			return c, e
		}

		_, err = strconv.Atoi(parts[2])
		if err != nil {
			e = fmt.Errorf("Destination port could not be converted to an integer")
			return c, e
		}
		dest := parts[1] + ":" + parts[2]

		c.Ports[port] = dest
	}
	fmt.Printf("Redir: %+v\n", redirFlags)

	// TLS PORT Parsing
	c.TLSPorts = make(intslice, 0)
	parts := strings.Split(tlsPorts, ",")
	if tlsPorts != "" {
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				e = fmt.Errorf("Could not parse a TLS port as an integer")
				break
			}
			c.TLSPorts = append(c.TLSPorts, v)
		}
	}
	return c, e
}
