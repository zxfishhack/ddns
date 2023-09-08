package namecheap

import "encoding/xml"

type Host struct {
	XMLName xml.Name `xml:"host"`
	HostId  string   `xml:"HostId,attr"`
	Name    string   `xml:"Name,attr"`
	Address string   `xml:"Address,attr"`
	Type    string   `xml:"Type,attr"`
	MXPref  string   `xml:"MXPref,attr"`
	TTL     string   `xml:"TTL,attr"`
}

type DomainDNSGetHostsResult struct {
	XMLName xml.Name `xml:"DomainDNSGetHostsResult"`
	Domain  string   `xml:"Domain,attr"`
	Hosts   []*Host  `xml:"host"`
}

type DomainDNSSetHostsResult struct {
	XMLName   xml.Name `xml:"DomainDNSSetHostsResult"`
	Domain    string   `xml:"Domain,attr"`
	IsSuccess bool     `xml:"IsSuccess,attr"`
}

type CommandResponse struct {
	XMLName                 xml.Name `xml:"CommandResponse"`
	DomainDNSGetHostsResult *DomainDNSGetHostsResult
	DomainDNSSetHostsResult *DomainDNSSetHostsResult
}

type ApiResponse struct {
	XMLName         xml.Name `xml:"ApiResponse"`
	CommandResponse CommandResponse
}

type AuthDetails struct {
	ParentUserType string `json:"ParentUserType"`
	ParentUserId   string `json:"ParentUserId"`
	UserId         string `json:"UserId"`
	UserName       string `json:"UserName"`
	ClientIp       string `json:"ClientIp"`
}

type setHostsRequest struct {
	AuthDetails AuthDetails `json:"authDetails"`
}
