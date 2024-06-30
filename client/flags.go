package client

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Opts struct {
	Cmd          string
	Format       string
	Free         bool
	Id           int
	Parent       int
	Pretty       bool
	SearchString string
}

func ParseCmdFlags(args []string) *Opts {
	params := &Opts{}
	params.Cmd = args[0]
	switch params.Cmd {
	case "ranges":
		flags := flag.NewFlagSet(params.Cmd, flag.ExitOnError)
		flags.BoolVar(&params.Pretty, "pretty", false, "pretty print range")
		flags.IntVar(&params.Id, "id", 0, "get range with given id")
		flags.BoolVar(&params.Free, "free", false, "get non-allocated IP addresses from given main range id")
		flags.IntVar(&params.Parent, "parent", 0, "get ranges with given parent")
		flags.StringVar(&params.SearchString, "search", "", "search string in range Name and CIDR fields")
		flags.StringVar(&params.Format, "format", "table", "format of the -free command")
		flags.StringVar(&params.SearchString, "s", "", "search string in range Name and CIDR fields")
		err := flags.Parse(args[1:])
		if err != nil {
			log.Printf("%s\n", err.Error())
		}

	case "domains":
		flags := flag.NewFlagSet(params.Cmd, flag.ExitOnError)
		flags.BoolVar(&params.Pretty, "pretty", false, "pretty print range")
		flags.IntVar(&params.Id, "id", 0, "get domain with id")
		flags.StringVar(&params.SearchString, "search", "", "search string in routing domain Id and VPCs fields")
		flags.StringVar(&params.SearchString, "s", "", "search string in routing domain Id and VPCs fields")
		err := flags.Parse(args[1:])
		if err != nil {
			log.Printf("%s\n", err.Error())
		}

	case "status":
		return &Opts{Cmd: "status"}

	default:
		fmt.Printf("unknown command: %s\n", params.Cmd)
		flag.Usage()
		os.Exit(1)
	}

	if os.Getenv("IPAM_PRETTY") == "true" {
		params.Pretty = true
	}

	return params
}

func Usage() {
	fmt.Printf(`Usage:
  ipam ranges                       : all subnet ranges
  ipam ranges -id <int>             : range with Subnet Id equal to <int>
  ipam ranges -id <int> -free       : number of free / non-allocated IP addresses for range
  ipam ranges -id <int> -free -format [ table (default) | json | number]
                                        - table  : displays CIDR, Allocated IPs, Available IPs and CIDR Name
                                        - json   : returns a JSON object
                                        - number : returns the number of available IP addresses  
  ipam ranges -parent <int>         : all subnet ranges with a parent Subnet Id matching <int>
  ipam ranges -s, -search <string>  : all subnet ranges that contains <string> in the subnet's Name or CIDR fields

  ipam domains                      : all routing domains
  ipam domains -id <int>            : routing domain with Id equal to <int>
  ipam domains -s, -search <string> : all routing domains that contains <string> in the domains's Name or VPCs fields

  ipam status                       : returns IPAM status

Global:
  -pretty                           : display indented JSON form of ranges/domains

Environment variables:
  IPAM_PRETTY=true                  : default to display indented JSON
  POKEMONIZE=true | ANONYMIZE=true  : anonymize range names and IPs

Help:
  ipam -h, -help                    : Print this help
  ipam ranges|domains -h, -help     : Print available flags for the command
  ipam -v | ipam -version           : Print IPAM Client version

  All available flags support two dashes; e.g. --id
`)
}
