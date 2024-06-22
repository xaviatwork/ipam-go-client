package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func parseCmdFlags(args []string) {
	cmd = args[0]
	switch cmd {
	case "ranges":
		flags := flag.NewFlagSet(cmd, flag.ExitOnError)
		flags.BoolVar(&pretty, "pretty", false, "pretty print range")
		flags.IntVar(&id, "id", 0, "get range with given id")
		flags.IntVar(&parent, "parent", 0, "get ranges with given parent")
		flags.StringVar(&searchString, "search", "", "search string in range Name and CIDR fields")
		err := flags.Parse(args[1:])
		if err != nil {
			log.Printf("%s\n", err.Error())
		}

	case "domains":
		flags := flag.NewFlagSet(cmd, flag.ExitOnError)
		flags.BoolVar(&pretty, "pretty", false, "pretty print range")
		flags.IntVar(&id, "id", 0, "get domain with id")
		flags.StringVar(&searchString, "search", "", "search string in routing domain Id and VPCs fields")
		err := flags.Parse(args[1:])
		if err != nil {
			log.Printf("%s\n", err.Error())
		}
	case "status":
		return
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`Usage:
  ipam ranges                   : all subnet ranges
  ipam ranges -id <int>         : range with Subnet Id equal to <int>
  ipam ranges -parent <int>     : all subnet ranges with a parent Subnet Id matching <int>
  ipam ranges -search <string>  : all subnet ranges that contains <string> in the subnet's Name or CIDR fields

  ipam domains                  : all routing domains
  ipam domains -id <int>        : routing domain with Id equal to <int>
  ipam domains -search <string> : all routing domains that contains <string> in the domains's Name or VPCs fields
`)
}
