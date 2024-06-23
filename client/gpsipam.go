package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/xaviatwork/ipam/ipamautopilot"
)

type GpsIpam struct {
	Source string
}

func (gpsipam GpsIpam) RangeById(id int) (*ipamautopilot.Range, error) {
	lzrange := &ipamautopilot.Range{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/ranges/%d", gpsipam.Source, id))
	if err != nil {
		return lzrange, err
	}
	if err := json.Unmarshal(b, &lzrange); err != nil {
		return lzrange, err
	}
	return lzrange, nil
}
func (gpsipam GpsIpam) Ranges() (*[]ipamautopilot.Range, error) {
	ranges := &[]ipamautopilot.Range{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/ranges", gpsipam.Source))
	if err != nil {
		return ranges, err
	}
	if err := json.Unmarshal(b, &ranges); err != nil {
		return ranges, err
	}
	return ranges, nil
}
func (gpsipam GpsIpam) RoutingDomainById(id int) (*ipamautopilot.RoutingDomain, error) {
	routingdomain := &ipamautopilot.RoutingDomain{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/domains/%d", gpsipam.Source, id))
	if err != nil {
		return routingdomain, err
	}
	if err := json.Unmarshal(b, &routingdomain); err != nil {
		return routingdomain, err
	}
	return routingdomain, nil
}
func (gpsipam GpsIpam) RoutingDomains() (*[]ipamautopilot.RoutingDomain, error) {
	domains := &[]ipamautopilot.RoutingDomain{}
	b, err := gpsipam.doRequest(fmt.Sprintf("%s/domains", gpsipam.Source))
	if err != nil {
		return domains, err
	}
	if err := json.Unmarshal(b, &domains); err != nil {
		return domains, err
	}
	return domains, nil
}

func (gpsipam GpsIpam) getToken() string {
	return os.Getenv("IPAM_TOKEN")
}

func (gpsipam GpsIpam) doRequest(url string) ([]byte, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("Authorization", "bearer "+gpsipam.getToken())

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

func (gpsipam GpsIpam) Status() error {
	b, err := gpsipam.doRequest(gpsipam.Source)
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
		if searchString(opts.SearchString, r.Name, r.Cidr) {
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
		if searchString(opts.SearchString, d.Name, d.Vpcs) {
			if opts.Pretty {
				fmt.Printf("%s\n", d.PrettyString())
				continue
			}
			fmt.Printf("%s", d.String())
		}
	}
}

// SearchString returns true if any of the ss[1:] strings contains ss[0]
func searchString(ss ...string) bool {
	searchString := strings.ToLower(ss[0])
	found := false
	for _, s := range ss[1:] {
		if strings.Contains(strings.ToLower(s), searchString) {
			found = true
		}
	}
	return found
}
