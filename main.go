package main

import (
	"fmt"
	"os"

	"github.com/xaviatwork/ipam-client/ipamclient"
)

func main() {
	ipam := ipamclient.LzIpam{BaseUrl: os.Getenv("IPAM_URL")}
	if ipam.BaseUrl == "" {
		fmt.Println("No IPAM url provided. Please set IPAM_URL env variable")
		os.Exit(1)
	}
	lzranges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM returned %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(lzranges)
<<<<<<< HEAD
=======

	lzrange, err := ipam.RangeById(2)
	if err != nil {
		fmt.Printf("IPAM returned %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("IPAM Range: %+v\n", lzrange)
>>>>>>> flags
}
