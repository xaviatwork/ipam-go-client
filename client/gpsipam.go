package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
