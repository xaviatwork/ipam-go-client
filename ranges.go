package main

import (
	"fmt"
	"os"

	"github.com/xaviatwork/ipam-client/ipamclient"
)

func getRangeById(id int, ipam ipamclient.LzIpam) {
	iprange, err := ipam.RangeById(id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s\n", iprange.String())
}

func getRangesWithSameParentSubnet(parent int, ipam ipamclient.LzIpam) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		if r.Parent_id == parent {
			fmt.Printf("%s", r.String())
		}
	}
}

func searchStringInRanges(searchString string, ipam ipamclient.LzIpam) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		if r.SearchString(searchString, r.Name, r.Cidr) {
			fmt.Printf("%s", r.String())
		}
	}
}
