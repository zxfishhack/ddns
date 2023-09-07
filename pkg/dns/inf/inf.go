package inf

import "errors"

type IDNSService interface {
	GetDomainRecord(subDomain, domainName string) (ip, rid string, err error)
	CreateDomainRecord(subDomain, domainName, ip string) (err error)
	UpdateDomainRecord(subDomain, domainName, ip, rid string) (err error)
}

var ErrNotAvailable = errors.New("dns service is not available")
