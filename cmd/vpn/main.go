package main

import (
	"flag"
	"fmt"

	"github.com/prophittcorey/vpn"
)

func main() {
	var vpns bool
	var ip string

	flag.BoolVar(&vpns, "subnets", false, "if specified, all VPN subnets will be dumped in one list")
	flag.StringVar(&ip, "check", "", "an ip address to analyze (returns information about the address)")

	flag.Parse()

	if vpns {
		for _, subnet := range vpn.Subnets() {
			fmt.Printf("%s\n", subnet.String())
		}

		return
	}

	if len(ip) > 0 {
		name, err := vpn.Check(ip)

		if err == nil {
			fmt.Printf("Looks like a '%s' address.\n", name)
		} else {
			fmt.Printf("Does not look like a vpn address.\n")
		}

		return
	}

	flag.Usage()
}
