package namecheap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zxfishhack/ddns/pkg/dns/inf"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	namecheapEndpoint = "https://api.namecheap.com/xml.response"
)

func init() {
	v := os.Getenv("NAMECHEAP_API")
	if v != "" {
		namecheapEndpoint = v
	}
}

type namecheapDns struct {
	apiKey   string
	userName string
	sourceIp string
}

func (d *namecheapDns) getHosts(domainName string) (list []*Host, err error) {
	q := url.Values{}
	d.addAuth(&q)
	q.Set("Command", "namecheap.domains.dns.getHosts")
	parts := strings.Split(domainName, ".")
	if len(parts) != 2 {
		err = errors.New("domainName is error")
		return
	}
	q.Set("SLD", parts[0])
	q.Set("TLD", parts[1])

	resp, err := http.Post(namecheapEndpoint, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code(%d) != 200", resp.StatusCode))
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var apiResponse ApiResponse
	err = xml.Unmarshal(b, &apiResponse)
	if err != nil {
		return
	}
	if apiResponse.CommandResponse.DomainDNSGetHostsResult != nil {
		list = apiResponse.CommandResponse.DomainDNSGetHostsResult.Hosts
	}
	return
}

func (d *namecheapDns) GetDomainRecord(subDomain, domainName string) (ip, rid string, err error) {
	list, err := d.getHosts(domainName)
	if err == nil {
		err = os.ErrNotExist
		for _, h := range list {
			if h.Name == subDomain {
				ip = h.Address
				rid = h.HostId
				err = nil
				break
			}
		}
	} else {
		err = os.ErrNotExist
	}

	q := url.Values{}
	d.addAuth(&q)
	q.Set("Command", "namecheap.domains.dns.getHosts")
	parts := strings.Split(domainName, ".")
	if len(parts) != 2 {
		err = errors.New("domainName is error")
		return
	}
	q.Set("SLD", parts[0])
	q.Set("TLD", parts[1])

	resp, err := http.Post(namecheapEndpoint, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code(%d) != 200", resp.StatusCode))
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var apiResponse ApiResponse
	err = xml.Unmarshal(b, &apiResponse)
	if err != nil {
		return
	}
	if apiResponse.CommandResponse.DomainDNSGetHostsResult != nil {
		list = apiResponse.CommandResponse.DomainDNSGetHostsResult.Hosts
	}

	return
}

func (d *namecheapDns) CreateDomainRecord(subDomain, domainName, ip string) (err error) {
	list, err := d.getHosts(domainName)
	if err != nil {
		return
	}
	q := url.Values{}
	d.addAuth(&q)
	q.Set("Command", "namecheap.domains.dns.setHosts")
	parts := strings.Split(domainName, ".")
	if len(parts) != 2 {
		err = errors.New("domainName is error")
		return
	}
	q.Set("SLD", parts[0])
	q.Set("TLD", parts[1])
	found := false
	for _, h := range list {
		if h.Name == subDomain && h.Type == "A" {
			h.Address = ip
			found = true
			break
		}
	}
	if !found {
		list = append(list, &Host{
			Name:    subDomain,
			Address: ip,
			Type:    "A",
			MXPref:  "10",
		})
	}

	for i, h := range list {
		q.Set("HostName"+strconv.Itoa(i+1), h.Name)
		q.Set("RecordType"+strconv.Itoa(i+1), h.Type)
		q.Set("Address"+strconv.Itoa(i+1), h.Address)
		q.Set("MXPref"+strconv.Itoa(i+1), h.MXPref)
		q.Set("TTL"+strconv.Itoa(i+1), h.TTL)
	}

	resp, err := http.Post(namecheapEndpoint, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code(%d) != 200", resp.StatusCode))
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var apiResponse ApiResponse
	err = xml.Unmarshal(b, &apiResponse)
	if err != nil {
		return
	}
	if apiResponse.CommandResponse.DomainDNSSetHostsResult == nil || !apiResponse.CommandResponse.DomainDNSSetHostsResult.IsSuccess {
		err = errors.New("setHosts failed")
	}
	return
}

func (d *namecheapDns) UpdateDomainRecord(subDomain, domainName, ip, rid string) (err error) {
	return d.CreateDomainRecord(subDomain, domainName, ip)
}

func (d *namecheapDns) addAuth(q *url.Values) {
	q.Set("ApiUser", d.userName)
	q.Set("ApiKey", d.apiKey)
	q.Set("UserName", d.userName)
	q.Set("ClientIp", d.sourceIp)
}

func New() (s inf.IDNSService, err error) {
	dc := &namecheapDns{
		apiKey:   os.Getenv("NAMECHEAP_API_KEY"),
		userName: os.Getenv("NAMECHEAP_USERNAME"),
		sourceIp: os.Getenv("NAMECHEAP_SOURCEIP"),
	}

	var u *url.URL
	u, err = url.Parse(namecheapEndpoint)
	if err != nil {
		return
	}
	q := url.Values{}
	dc.addAuth(&q)
	q.Set("Command", "namecheap.domains.getList")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = inf.ErrNotAvailable
		return
	}

	s = dc

	return
}
