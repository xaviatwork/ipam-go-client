package ipamclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type LzIpam struct {
	BaseUrl string
}

func (lzipam *LzIpam) getToken() string {
	return os.Getenv("IPAM_TOKEN")
}

func (lzipam *LzIpam) Ranges() (*[]Range, error) {
	ranges := &[]Range{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/ranges", lzipam.BaseUrl))
	if err != nil {
		return ranges, err
	}
	if err := json.Unmarshal(b, &ranges); err != nil {
		return ranges, err
	}
	return ranges, nil
}

func (lzipam *LzIpam) RangeById(id int) (*Range, error) {
	lzrange := &Range{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/ranges/%d", lzipam.BaseUrl, id))
	if err != nil {
		return lzrange, err
	}
	if err := json.Unmarshal(b, &lzrange); err != nil {
		return lzrange, err
	}
	return lzrange, nil
}

func (lzipam *LzIpam) RoutingDomains() (*[]RoutingDomain, error) {
	domains := &[]RoutingDomain{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/domains", lzipam.BaseUrl))
	if err != nil {
		return domains, err
	}
	if err := json.Unmarshal(b, &domains); err != nil {
		return domains, err
	}
	return domains, nil
}

func (lzipam *LzIpam) RoutingDomainById(id int) (*RoutingDomain, error) {
	routingdomain := &RoutingDomain{}
	b, err := lzipam.doRequest(fmt.Sprintf("%s/domains/%d", lzipam.BaseUrl, id))
	if err != nil {
		return routingdomain, err
	}
	if err := json.Unmarshal(b, &routingdomain); err != nil {
		return routingdomain, err
	}
	return routingdomain, nil
}

func (lzipam *LzIpam) doRequest(url string) ([]byte, error) {
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

func (lzipam *LzIpam) Status() error {
	b, err := lzipam.doRequest(lzipam.BaseUrl)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
