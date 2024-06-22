package ipamclient

type RoutingDomain struct {
	Id   int    `db:"routing_domain_id"`
	Name string `db:"name"`
	Vpcs string `db:"vpcs"` // associated VPCs that should be tracked for subnet creation
}

type Range struct {
	Subnet_id         int    `db:"subnet_id"`
	Parent_id         int    `db:"parent_id"`
	Routing_domain_id int    `db:"routing_domain_id"`
	Name              string `db:"name"`
	Cidr              string `db:"cidr"`
}

type IpamAutopilot interface {
	Ranges() (*[]Range, error)
	RangeById(id int) (*Range, error)
	// RoutingDomains() ([]RoutingDomain, error)
	// RoutingDomainById(id int) (*RoutingDomain, error)
}
