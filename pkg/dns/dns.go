package dns

import (
	"github.com/zxfishhack/ddns/pkg/dns/aliyun"
	"github.com/zxfishhack/ddns/pkg/dns/inf"
	"github.com/zxfishhack/ddns/pkg/dns/namecheap"
)

var svc inf.IDNSService

func Init() (err error) {
	// try init aliyun
	svc, err = aliyun.New()
	if err == nil {
		return
	}
	// try init namecheap
	svc, err = namecheap.New()
	if err == nil {
		return
	}
	return
}

func GetDomainRecord(subDomain, domainName string) (ip, rid string, err error) {
	err = inf.ErrNotAvailable
	if svc != nil {
		return svc.GetDomainRecord(subDomain, domainName)
	}
	return
}

func CreateDomainRecord(subDomain, domainName, ip string) (err error) {
	err = inf.ErrNotAvailable
	if svc != nil {
		return svc.CreateDomainRecord(subDomain, domainName, ip)
	}
	return
}

func UpdateDomainRecord(subDomain, domainName, ip, rid string) (err error) {
	err = inf.ErrNotAvailable
	if svc != nil {
		return svc.UpdateDomainRecord(subDomain, domainName, ip, rid)
	}
	return
}
