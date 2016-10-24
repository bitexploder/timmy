package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Ports map[int]string
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

func parseFlags() Config {
	//var redir = flag.String("redir", "", "Redirect a port to a destination [orig_port]:[dest_server]:[dest_port]")
	var redirFlags redir

	flag.Var(&redirFlags, "redir", "redirect port to target server:port port:serverip:serverport")

	flag.Parse()

	fmt.Printf("Redir: %+v\n", redirFlags)

	c := Config{
		Ports: make(map[int]string),
	}

	return c
}
