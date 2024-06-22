package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xaviatwork/ipam-client/ipamclient"
)

var cmd string
var id int
var parent int
var searchString string

func main() {
	// main flags (just -h, --help)
	flag.Usage = func() { usage() }
	flag.Parse()
	// cmd flags
	parseCmdFlags(flag.Args())

	ipam := ipamclient.LzIpam{BaseUrl: os.Getenv("IPAM_URL")}
	if ipam.BaseUrl == "" {
		fmt.Println("No IPAM url provided. Please set IPAM_URL env variable")
		os.Exit(1)
	}

	switch cmd {
	case "ranges":
		switch {
		case id != 0:
			// if the range is not found, IPAM Autopilot returns a 503 Service Unavailable error
			// https://github.com/GoogleCloudPlatform/professional-services/blob/main/tools/ipam-autopilot/container/api.go#L81
			getRangeById(id, ipam)
		case parent != 0:
			getRangesWithParent(parent, ipam)
		case searchString != "":
			searchStringInRanges(searchString, ipam)
		default: // all ranges
			searchStringInRanges("", ipam)
		}
	case "domains":
		switch {
		case id != 0:
			// if the range is not found, IPAM Autopilot returns a 503 Service Unavailable error
			// https://github.com/GoogleCloudPlatform/professional-services/blob/main/tools/ipam-autopilot/container/api.go#L81
			getDomainById(id, ipam)
		case searchString != "":
			searchStringInDomains(searchString, ipam)
		default: // all Domainss
			searchStringInDomains("", ipam)
		}
	case "status":
		ipam.Status()
	}
}
