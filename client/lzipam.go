package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/xaviatwork/ipam/ipamautopilot"
)

type LzIpam struct {
	Source string
}

func (lzipam LzIpam) RangeById(id int) (*ipamautopilot.Range, error) {
	lzrange := &ipamautopilot.Range{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/ranges/%d", lzipam.Source, id))
	if err != nil {
		return lzrange, err
	}
	if err := json.Unmarshal(b, &lzrange); err != nil {
		return lzrange, err
	}
	return lzrange, nil
}
func (lzipam LzIpam) Ranges() (*[]ipamautopilot.Range, error) {
	ranges := &[]ipamautopilot.Range{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/ranges", lzipam.Source))
	if err != nil {
		return ranges, err
	}
	if err := json.Unmarshal(b, &ranges); err != nil {
		return ranges, err
	}
	return ranges, nil
}
func (lzipam LzIpam) RoutingDomainById(id int) (*ipamautopilot.RoutingDomain, error) {
	routingdomain := &ipamautopilot.RoutingDomain{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/domains/%d", lzipam.Source, id))
	if err != nil {
		return routingdomain, err
	}
	if err := json.Unmarshal(b, &routingdomain); err != nil {
		return routingdomain, err
	}
	return routingdomain, nil
}
func (lzipam LzIpam) RoutingDomains() (*[]ipamautopilot.RoutingDomain, error) {
	domains := &[]ipamautopilot.RoutingDomain{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/domains", lzipam.Source))
	if err != nil {
		return domains, err
	}
	if err := json.Unmarshal(b, &domains); err != nil {
		return domains, err
	}
	return domains, nil
}

// func (lzipam LzIpam) Source() string {
// 	return os.Getenv("IPAM_SOURCE")
// }

func (lzipam LzIpam) getToken() string {
	return os.Getenv("IPAM_TOKEN")
}

func (lzipam LzIpam) doRequest(url string) ([]byte, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("Authorization", "bearer "+lzipam.getToken())

	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return []byte{}, fmt.Errorf("http error %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (lzipam LzIpam) Status() error {
	b, err := lzipam.doRequest(lzipam.Source)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func GetRangeById(ipam ipamautopilot.Ipam, opts Opts) {
	iprange, err := ipam.RangeById(opts.Id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	if opts.Pretty {
		fmt.Printf("%s\n", iprange.PrettyString())
		return
	}
	fmt.Printf("%s\n", iprange.String())
}

func GetRangesWithParent(ipam ipamautopilot.Ipam, opts Opts) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		if r.Parent_id == opts.Parent {
			if opts.Pretty {
				fmt.Printf("%s\n", r.PrettyString())
				continue
			}
			fmt.Printf("%s", r.String())
		}
	}
}

func SearchStringInRanges(ipam ipamautopilot.Ipam, opts Opts) {
	ipRanges, err := ipam.Ranges()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, r := range *ipRanges {
		if ipamautopilot.SearchString(opts.SearchString, r.Name, r.Cidr) {
			if opts.Pretty {
				fmt.Printf("%s\n", r.PrettyString())
				continue
			}
			fmt.Printf("%s", r.String())
		}
	}
}

func GetDomainById(ipam ipamautopilot.Ipam, opts Opts) {
	domain, err := ipam.RoutingDomainById(opts.Id)
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	if opts.Pretty {
		fmt.Printf("%s\n", domain.PrettyString())
		return
	}
	fmt.Printf("%s\n", domain.String())
}
func SearchStringInDomains(ipam ipamautopilot.Ipam, opts Opts) {
	domains, err := ipam.RoutingDomains()
	if err != nil {
		fmt.Printf("IPAM response: %s\n", err.Error())
		os.Exit(1)
	}
	for _, d := range *domains {
		if ipamautopilot.SearchString(opts.SearchString, d.Name, d.Vpcs) {
			if opts.Pretty {
				fmt.Printf("%s\n", d.PrettyString())
				continue
			}
			fmt.Printf("%s", d.String())
		}
	}
}
