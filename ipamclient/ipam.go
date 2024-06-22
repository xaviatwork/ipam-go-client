package ipamclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type RoutingDomain struct {
	Id   int    `json:"id"`
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

func (r *Range) PrettyString() string {
	pretty := bytes.Buffer{}
	b, err := json.Marshal(r)
	if err != nil {
		return r.String()
	}
	if err := json.Indent(&pretty, b, "", "  "); err != nil {
		return r.String()
	}
	return pretty.String()
}

func (d *RoutingDomain) String() string {
	return fmt.Sprintf(`{"id": %d, "name": "%s", "vpcs": "%s"}`, d.Id, d.Name, d.Vpcs)
}

func (d *RoutingDomain) PrettyString() string {
	pretty := bytes.Buffer{}
	b, err := json.Marshal(d)
	if err != nil {
		return d.String()
	}
	if err := json.Indent(&pretty, b, "", "  "); err != nil {
		return d.String()
	}
	return pretty.String()
}

// SearchString returns true if any of the ss[1:] strings contains ss[0]
//
//	All strings are converted to lowercase to compare them.
func SearchString(ss ...string) bool {
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
	RoutingDomains() ([]RoutingDomain, error)
	RoutingDomainById(id int) (*RoutingDomain, error)
}
