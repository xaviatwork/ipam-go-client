package ipamclient

import (
	"fmt"
	"strings"
)

type RoutingDomain struct {
	Id   int    `json:"domain"`
	Name string `json:"name"`
	Vpcs string `json:"vpcs"` // associated VPCs that should be tracked for subnet creation
}

type Range struct {
	Cidr              string `json:"cidr"`
	Name              string `json:"name"`
	Parent_id         int    `json:"parent"`
	Routing_domain_id int    `json:"domain"`
	Subnet_id         int    `json:"id"`
}

func (r *Range) String() string {
	return fmt.Sprintf(`{"cidr": "%s", "name": "%s", "parent": %d, "domain": %d, "id": %d}`, r.Cidr, r.Name, r.Parent_id, r.Routing_domain_id, r.Subnet_id)
}

// SearchString returns true if any of the ss[1:] strings contains ss[0]
//
//	All strings are converted to lowercase to compare them.
func (r *Range) SearchString(ss ...string) bool {
	searchString := strings.ToLower(ss[0])
	found := false
	for _, s := range ss[1:] {
		if strings.Contains(strings.ToLower(s), searchString) {
			found = true
		}
	}
	return found
}

type IpamAutopilot interface {
	Ranges() (*[]Range, error)
	RangeById(id int) (*Range, error)
	// RoutingDomains() ([]RoutingDomain, error)
	// RoutingDomainById(id int) (*RoutingDomain, error)
}
