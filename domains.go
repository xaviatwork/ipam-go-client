package main

import (
	"fmt"
	"os"

	"github.com/xaviatwork/ipam-client/ipamclient"
)

func getDomainById(id int, ipam ipamclient.LzIpam) {
	domain, err := ipam.RoutingDomainById(id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	if pretty {
		fmt.Printf("%s\n", domain.PrettyString())
		return
	}
	fmt.Printf("%s\n", domain.String())
}
func searchStringInDomains(searchString string, ipam ipamclient.LzIpam) {
	domains, err := ipam.RoutingDomains()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, d := range *domains {
		if ipamclient.SearchString(searchString, d.Name, d.Vpcs) {
			if pretty {
				fmt.Printf("%s\n", d.PrettyString())
				continue
			}
			fmt.Printf("%s", d.String())
		}
	}
}
