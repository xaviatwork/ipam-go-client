package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xaviatwork/ipam/client"
)

func main() {
	// main flags (just -h, --help)
	flag.Usage = func() { client.Usage() }
	flag.Parse()
	// cmd flags
	if len(os.Args) < 2 {
		client.Usage()
		os.Exit(1)
	}
	opts := client.ParseCmdFlags(flag.Args())

	ipam := client.GpsIpam{} // Implements ipamautopilot.Ipam
	ipam.Source = os.Getenv("IPAM_SOURCE")
	if ipam.Source == "" {
		fmt.Println("No IPAM source provided. Please set IPAM_SOURCE env variable")
		os.Exit(1)
	}

	switch opts.Cmd {
	case "ranges":
		switch {
		case opts.Id != 0:
			// if the range is not found, IPAM Autopilot returns a 503 Service Unavailable error
			// https://github.com/GoogleCloudPlatform/professional-services/blob/main/tools/ipam-autopilot/container/api.go#L81
			if opts.Free {
				client.GetNonAllocatedIPs(ipam, *opts)
				os.Exit(0)
			}
			client.GetRangeById(ipam, *opts)
		case opts.Parent != 0:
			client.GetRangesWithParent(ipam, *opts)
		case opts.SearchString != "":
			client.SearchStringInRanges(ipam, *opts)
		default: // all ranges
			opts.SearchString = ""
			client.SearchStringInRanges(ipam, *opts)
		}
	case "domains":
		switch {
		case opts.Id != 0:
			// if the domain is not found, IPAM Autopilot returns a 503 Service Unavailable error
			// https://github.com/GoogleCloudPlatform/professional-services/blob/main/tools/ipam-autopilot/container/api.go#L370
			client.GetDomainById(ipam, *opts)
		case opts.SearchString != "":
			client.SearchStringInDomains(ipam, *opts)
		default: // all Domainss
			opts.SearchString = ""
			client.SearchStringInDomains(ipam, *opts)
		}
	case "status":
		log.Println("hey there")
		ipam.Status()
	}
}
